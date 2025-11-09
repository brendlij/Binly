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
docker-compose up
```

Visit `http://localhost` in your browser.

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
docker-compose up
```

Runs on `http://localhost`

## License

MIT

---

Made with minimal aesthetics in mind. [Source](https://github.com/brendlij/Binly)
