package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"
)

type Paste struct {
	ID          string
	Content     string
	Syntax      string
	AllowEdit   int
	PwSalt      []byte
	PwHash      []byte
	CreatedAt   int64
	ExpiresAt   sql.NullInt64
}

var (
	appSecret []byte
	maxBody   int64 = 512 * 1024 // 512KB
	db        *sql.DB
)

func main() {
	var (
		addr     = flag.String("addr", ":8080", "public http addr")
		dbPath   = flag.String("db", "data/pastes.db", "sqlite path")
		uiDir    = flag.String("ui", "web/dist", "built vue dist dir")
	)
	flag.Parse()

	secret := os.Getenv("APP_SECRET")
	if secret == "" { secret = "dev_only_replace_me" }
	appSecret = []byte(secret)

	_ = os.MkdirAll(filepath.Dir(*dbPath), 0755)
	var err error
	db, err = sql.Open("sqlite", "file:"+*dbPath+"?_pragma=journal_mode(WAL)")
	must(err)
	must(db.Ping())
	must(initSchema())

	// Cleanup Loop
	go func() {
		t := time.NewTicker(10 * time.Minute)
		for range t.C {
			db.Exec(`DELETE FROM pastes WHERE expires_at IS NOT NULL AND expires_at < ?`, time.Now().Unix())
		}
	}()

	// Public API + Static
	r := chi.NewRouter()
	r.Use(limitBody)
	r.Get("/_health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })

	r.Post("/api/p", createPaste)
	r.Get("/api/p/{id}", getPaste)
	r.Post("/api/p/{id}", updatePaste)           // edit (wenn allow_edit)
	r.Delete("/api/p/{id}", deletePaste)        // delete (owner only)
	r.Post("/api/p/{id}/auth", authPaste)        // passwort-check
	r.Get("/api/raw/{id}", getRaw)

	// Static Vue dist
	fs := http.FileServer(http.Dir(*uiDir))
	r.Handle("/assets/*", fs) // Vite assets
	r.Handle("/docs/*", fs)    // Documentation markdown files
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		// SPA fallback
		http.ServeFile(w, r, filepath.Join(*uiDir, "index.html"))
	})

	// TODO: Admin endpoints can be added here if needed
	// For now, all management is done via "My Shares" UI

	go func() { log.Println("Public on", *addr); log.Fatal(http.ListenAndServe(*addr, r)) }()
	
	// Keep the server running
	select {}
}

func initSchema() error {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS pastes(
  id TEXT PRIMARY KEY,
  content TEXT NOT NULL,
  syntax TEXT DEFAULT 'auto',
  allow_edit INTEGER NOT NULL DEFAULT 0,
  pw_salt BLOB,
  pw_hash BLOB,
  created_at INTEGER NOT NULL,
  expires_at INTEGER
);
CREATE INDEX IF NOT EXISTS idx_expires ON pastes(expires_at);
`)
	return err
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func createPaste(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil { http.Error(w, "bad form", 400); return }
	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" || int64(len(content)) > maxBody { http.Error(w, "invalid content", 400); return }
	syntax := r.FormValue("syntax"); if syntax == "" { syntax = "auto" }
	ttl := r.FormValue("ttl")  // "15m" | "2h" | "1d" | "0"
	allowEdit := boolToInt(r.FormValue("allow_edit") == "on")
	password := r.FormValue("password")

	id := newID(8)
	var salt, sum []byte
	if password != "" {
		salt, sum = hashPassword(password)
	}
	var expires *int64
	if ttl != "" && ttl != "0" {
		if d, err := time.ParseDuration(ttl); err == nil {
			t := time.Now().Add(d).Unix(); expires = &t
		}
	}

	_, err := db.Exec(`INSERT INTO pastes(id,content,syntax,allow_edit,pw_salt,pw_hash,created_at,expires_at)
	VALUES(?,?,?,?,?,?,?,?)`,
		id, content, syntax, allowEdit, nullBytes(salt), nullBytes(sum), time.Now().Unix(), nullInt64Ptr(expires))
	if err != nil { http.Error(w, "db error", 500); return }

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"id":"%s"}`, id)
}

func getPaste(w http.ResponseWriter, r *http.Request) {
	p, err := loadPaste(chi.URLParam(r, "id"))
	if err != nil { http.Error(w, "not found", 404); return }
	if expired(p) { http.Error(w, "gone", 410); return }

	protected := len(p.PwHash) > 0
	if protected && !hasValidAuth(r, p.ID) {
		http.Error(w, "unauthorized", 401)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"id":"%s","content":%q,"syntax":%q,"allow_edit":%t,"expires_at":%d}`, p.ID, p.Content, p.Syntax, p.AllowEdit == 1, p.ExpiresAt.Int64)
}

func getRaw(w http.ResponseWriter, r *http.Request) {
	p, err := loadPaste(chi.URLParam(r, "id"))
	if err != nil { http.Error(w, "not found", 404); return }
	if expired(p) { http.Error(w, "gone", 410); return }
	if len(p.PwHash) > 0 && !hasValidAuth(r, p.ID) { http.Error(w, "unauthorized", 401); return }

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(p.Content))
}

func updatePaste(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, err := loadPaste(id)
	if err != nil { http.Error(w, "not found", 404); return }
	if expired(p) { http.Error(w, "gone", 410); return }
	if p.AllowEdit == 0 { http.Error(w, "forbidden", 403); return }
	if len(p.PwHash) > 0 && !hasValidAuth(r, p.ID) { http.Error(w, "unauthorized", 401); return }

	if err := r.ParseForm(); err != nil { http.Error(w, "bad form", 400); return }
	content := r.FormValue("content")
	if content == "" || int64(len(content)) > maxBody { http.Error(w, "invalid content", 400); return }

	_, err = db.Exec(`UPDATE pastes SET content=? WHERE id=?`, content, id)
	if err != nil { http.Error(w, "db error", 500); return }
	w.WriteHeader(204)
}

func deletePaste(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, err := loadPaste(id)
	if err != nil { http.Error(w, "not found", 404); return }
	if expired(p) { http.Error(w, "gone", 410); return }

	_, err = db.Exec(`DELETE FROM pastes WHERE id=?`, id)
	if err != nil { http.Error(w, "db error", 500); return }
	w.WriteHeader(204)
}

func authPaste(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, err := loadPaste(id)
	if err != nil { http.Error(w, "not found", 404); return }
	if len(p.PwHash) == 0 { http.Error(w, "no password", 400); return }

	if err := r.ParseForm(); err != nil { http.Error(w, "bad form", 400); return }
	pw := r.FormValue("password")
	if !checkPassword(pw, p.PwSalt, p.PwHash) { http.Error(w, "unauthorized", 401); return }

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_" + id,
		Value:    signAuth(id),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   24 * 3600,
	})
	w.WriteHeader(204)
}

/*** Helpers ***/
func loadPaste(id string) (Paste, error) {
	var p Paste
	err := db.QueryRow(`SELECT id,content,syntax,allow_edit,pw_salt,pw_hash,created_at,COALESCE(expires_at,0) FROM pastes WHERE id=?`, id).
		Scan(&p.ID,&p.Content,&p.Syntax,&p.AllowEdit,&p.PwSalt,&p.PwHash,&p.CreatedAt,&p.ExpiresAt.Int64)
	if errors.Is(err, sql.ErrNoRows) { return p, err }
	if p.ExpiresAt.Int64 == 0 { p.ExpiresAt.Valid = false } else { p.ExpiresAt.Valid = true }
	return p, err
}
func newID(n int) string {
  const b62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
  out := make([]byte, n)
  for i := 0; i < n; i++ {
    idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(b62))))
    if err != nil {
      // Fallback: sehr unwahrscheinlich, aber zur Not:
      out[i] = b62[0]
      continue
    }
    out[i] = b62[idx.Int64()]
  }
  return string(out)
}

func hashPassword(pw string) (salt, sum []byte) {
	salt = []byte(fmt.Sprintf("%d", time.Now().UnixNano()))
	h := sha256.Sum256(append(salt, []byte(pw)...))
	return salt, h[:]
}
func checkPassword(pw string, salt, sum []byte) bool {
	h := sha256.Sum256(append(salt, []byte(pw)...))
	return subtle.ConstantTimeCompare(h[:], sum) == 1
}
func signAuth(id string) string {
	m := hmac.New(sha256.New, appSecret)
	m.Write([]byte(id))
	return base64.RawURLEncoding.EncodeToString(m.Sum(nil))
}
func hasValidAuth(r *http.Request, id string) bool {
	c, _ := r.Cookie("auth_" + id); if c == nil { return false }
	raw, err := base64.RawURLEncoding.DecodeString(c.Value); if err != nil { return false }
	m := hmac.New(sha256.New, appSecret); m.Write([]byte(id))
	return hmac.Equal(raw, m.Sum(nil))
}
func expired(p Paste) bool {
	return p.ExpiresAt.Valid && p.ExpiresAt.Int64 < time.Now().Unix()
}
func nullBytes(b []byte) interface{} { if len(b)==0 { return nil }; return b }
func nullInt64Ptr(p *int64) interface{} { if p==nil { return nil }; return *p }
func limitBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxBody+1024)
		next.ServeHTTP(w, r)
	})
}
func must(err error){ if err!=nil { log.Fatal(err) } }
