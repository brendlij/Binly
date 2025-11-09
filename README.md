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

# Start with Docker Compose
docker-compose up
```

Visit `http://localhost` in your browser.

## Architecture

- **Backend**: Go (REST API)
- **Frontend**: Vue 3 + TypeScript + Vite
- **Syntax Highlighting**: highlight.js
- **Icons**: Tabler Icons
- **Styling**: CSS Variables, minimal design system

## Development

```bash
# Frontend
cd web
bun dev

# Backend
cd backend
go run ./cmd/server/main.go
```

## License

MIT

---

Made with minimal aesthetics in mind. [Source](https://github.com/brendlij/Binly)
