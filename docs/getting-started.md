# Getting Started

## Installation

Clone the repository and start with Docker:

```bash
git clone https://github.com/brendlij/Binly.git
cd Binly
export APP_SECRET=$(openssl rand -hex 32)
docker compose up -d
```

## Access

- **Web UI**: `http://localhost:7700`
- **In-App Docs**: `http://localhost:7700/docs`
- **My Shares**: `http://localhost:7700/my-shares`

## Create Your First Paste

1. Go to `http://localhost:7700`
2. Write or paste your code
3. Select a syntax (optional)
4. Click **Create**
5. Share the link!

## Track Your Shares

View all your shared pastes in **"My Shares"** at `http://localhost:7700/my-shares`

- 📋 See all your pastes at a glance
- ⏱️ Live TTL countdown
- 🎯 Quick copy/view/delete actions

## Optional Settings

- **TTL (Time to Live)**: Auto-delete after 15 minutes, 1 hour, 1 day, 7 days, or never
- **Allow Editing**: Others can edit this paste
- **Password**: Protect with a password
- **Syntax**: Auto-detect or manually select language
