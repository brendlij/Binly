# Binly

Share code and text instantly. No signup, no clutter, no compromise.

![Binly Screenshot](images/mainscreen-dark.png)

## Quick Start

```bash
git clone https://github.com/brendlij/Binly.git
cd Binly
export APP_SECRET=$(openssl rand -hex 32)
docker compose up -d
```

Visit `http://localhost:7700`

## Features

- 📋 Clean code sharing with syntax highlighting
- ⏱️ Temporary expiration (15m, 1h, 1d, 7d, or never)
- ✏️ Optional editing for shared pastes
- 🔒 Password protection
- 🌓 Dark/Light mode
- 📱 Fully responsive
- 🎯 "My Shares" to track your pastes

## Documentation

- [Getting Started](docs/getting-started.md) — Setup and configuration
- [API Reference](docs/api.md) — Complete endpoint documentation
- [User Guide](docs/guide.md) — Features and usage

Visit `http://localhost:7700/docs` for in-app documentation.

## Tech Stack

- **Backend**: Go 1.25 with chi router
- **Frontend**: Vue 3 + TypeScript + Vite
- **Database**: SQLite with WAL mode
- **Styling**: Syntax highlighting via highlight.js
- **Deployment**: Docker multi-stage build

## Development

```bash
# Frontend (dev mode)
cd web && bun install && bun dev

# Backend (dev mode)
cd backend/cmd/server && go run . --db ../../data/pastes.db --ui ../../web/dist --addr :8080
```

Then:

- Frontend: `http://localhost:5173` (Vite)
- Backend: `http://localhost:8080`
- Full stack with Docker: `docker compose up -d` → `http://localhost:7700`

## API Endpoints

### Public Endpoints

- `POST /api/p` — Create paste
- `GET /api/p/{id}` — Get paste
- `POST /api/p/{id}` — Update paste (if editable)
- `POST /api/p/{id}/auth` — Authenticate with password
- `GET /api/raw/{id}` — Get raw plain text

See [API Reference](docs/api.md) for details.

## License

GNU General Public License v3 — See [LICENSE](LICENSE)
