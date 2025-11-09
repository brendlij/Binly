# Binly

Minimal, fast code sharing. No clutter.

## Features

- **Minimal Design**: Clean, distraction-free interface
- **Dark & Light Mode**: Seamless theme switching with persistence
- **Fast & Lightweight**: Built with Go backend, Vue 3 frontend
- **Simple API**: Easy to integrate and use
- **Anonymous**: No accounts required

## Quick Start

```bash
# Clone the repository
git clone https://github.com/brendlij/Binly.git
cd Binly

# Generate a random secret for production
export APP_SECRET=$(openssl rand -hex 32)

# Start with Docker Compose
docker compose up -d
```

Visit `http://localhost:7700` in your browser.

## API Endpoints

### Public API

- `POST /api/p` - Create paste
- `GET /api/p/{id}` - Get paste
- `POST /api/p/{id}` - Edit paste (if allowed)
- `POST /api/p/{id}/auth` - Authenticate with password
- `GET /api/raw/{id}` - Get raw paste content

### Admin API (Port 7701)

> ⚠️ No authentication - only expose internally or use Cloudflare Tunnel with authentication

- `GET /admin` - List all pastes
- `POST /admin/delete/{id}` - Delete specific paste
- `POST /admin/purge` - Delete all expired pastes

## Ports

| Service   | Port | Access                                               |
| --------- | ---- | ---------------------------------------------------- |
| Web UI    | 7700 | `http://localhost:7700` or `http://192.168.x.x:7700` |
| Admin API | 7701 | `http://localhost:7701` or `http://192.168.x.x:7701` |

Local network only by default. Use Cloudflare Tunnel or reverse proxy for external access.

## Architecture

- **Backend**: Go (REST API)
- **Frontend**: Vue 3 + TypeScript + Vite
- **Syntax Highlighting**: highlight.js
- **Icons**: Tabler Icons via Iconify
- **Styling**: CSS Variables, minimal design system

## Development

### Local Setup

```bash
# Frontend (requires bun)
cd web
bun install
bun dev

# Backend (requires Go 1.23+)
cd backend/cmd/server
go run . --db ../../data/pastes.db --ui ../../web/dist --addr :8080 --admin :8787
```

Frontend runs on `http://localhost:5173`, backend serves on `http://localhost:8080`

### Production

```bash
export APP_SECRET=$(openssl rand -hex 32)
docker compose up -d
```

- Web UI: `http://localhost:7700`
- Admin API: `http://localhost:7701`

## License

MIT

---

Made with minimal aesthetics in mind. [Source](https://github.com/brendlij/Binly)
