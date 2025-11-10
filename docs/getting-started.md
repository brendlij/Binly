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
- **Admin API**: `http://localhost:7701`

## Create Your First Paste

1. Go to `http://localhost:7700`
2. Write or paste your code
3. Select a syntax (optional)
4. Click **Create**
5. Share the link!

## Optional Settings

- **TTL (Time to Live)**: Auto-delete after 15 minutes, 2 hours, 1 day, or never
- **Allow Editing**: Others can edit this paste
- **Password**: Protect with a password
