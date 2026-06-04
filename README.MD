# 🟢 ping-uptime

A self-hosted **uptime monitoring** service built with Go (Echo v4) and Vue 3. Monitor your websites, APIs, and services — get instant **Discord notifications** when something goes down or comes back up.

---

## ✨ Features

- 🔍 **HTTP/HTTPS uptime monitoring** — ping your endpoints at configurable intervals
- 🔔 **Discord webhook notifications** — instant alerts on status change (down / recovered)
- 📊 **Dashboard** — view real-time status, response time, and uptime history
- 🔐 **JWT authentication** — secure API and dashboard access
- 🗄️ **Multi-database support** — MySQL or PostgreSQL via GORM
- 🧩 **Modular architecture** — Domain-Driven Design, easy to extend
- ⚡ **Single binary** — Go backend serves the Vue frontend from `static/`

---

## 🛠️ Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go 1.25, Echo v4 |
| Frontend | Vue 3, Vue Router, Vite |
| ORM | GORM |
| Auth | JWT (`golang-jwt`) |
| Notifications | Discord Webhook |
| Database | MySQL / PostgreSQL |
| Logger | Uber Zap |
| Config | godotenv (`.env`) |

---

## 🚀 Getting Started

### Prerequisites

- Go `>= 1.21`
- Node.js `>= 20`
- MySQL or PostgreSQL

### 1. Clone the repo

```bash
git clone https://github.com/terarush/ping-uptime.git
cd ping-uptime
```

### 2. Configure environment

```bash
cp .env.example .env
```

Edit `.env` with your settings:

```env
APP_NAME=ping-uptime
SERVER_MODE=info
PORT=8080
HTTP_TIMEOUT=60

# Database
DB_DRIVER=mysql         # mysql | postgres
DB_HOST=localhost
DB_PORT=3306
DB_NAME=ping_uptime
DB_USERNAME=root
DB_PASSWORD=your_password

# Connection Pool
POOL_CONN_IDLE=200
POOL_CONN_MAX=300
POOL_CONN_LIFETIME=60

# JWT
JWT_SIGNATURE_KEY=your_secret_key
JWT_DAY_EXPIRED=60

# Discord
DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/your/webhook
```

### 3. Run the backend

```bash
go run main.go
# API available at http://localhost:8080/api
```

### 4. Run the frontend (development)

```bash
cd web
npm install
npm run dev
# Frontend at http://localhost:5173
# /api requests are proxied to :8080 automatically
```

---

## 🏗️ Project Structure

```
ping-uptime/
├── main.go                     # Entry point
├── .env.example                # Environment variable template
├── static/                     # Vue production build output
├── web/                        # Vue 3 frontend
│   ├── src/
│   ├── vite.config.ts          # Dev proxy + build config
│   └── package.json
└── internal/
    ├── app/                    # Application bootstrap
    │   ├── app.go              # Router, modules, SPA handler
    │   └── module.go           # Module interface
    └── pkg/
        ├── config/             # .env loader (godotenv)
        ├── database/           # GORM setup
        ├── jwt/                # JWT service
        ├── logger/             # Zap logger
        ├── middleware/         # Auth middleware
        ├── server/             # HTTP server + graceful shutdown
        └── utils/              # Response helpers
└── modules/
    ├── auth/                   # Login & register
    └── users/                  # User CRUD
```

---

## 🔔 Discord Notifications

When a monitored endpoint goes **down** or **recovers**, ping-uptime sends a Discord embed message via webhook:

**Down alert:**
```
🔴 [DOWN] https://yoursite.com
Status   : 503 Service Unavailable
Checked  : 2026-06-04 10:00:00 UTC
Response : timeout after 30s
```

**Recovery alert:**
```
🟢 [UP] https://yoursite.com
Status   : 200 OK
Checked  : 2026-06-04 10:05:00 UTC
Response : 142ms
```

### Setup Discord Webhook

1. Open your Discord server → **Settings → Integrations → Webhooks**
2. Click **New Webhook**, copy the URL
3. Paste into `.env`:
   ```env
   DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/...
   ```

---

## 🌐 API Endpoints

All endpoints are prefixed with `/api`.

### Auth
| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/auth/register` | Register new user |
| `POST` | `/api/auth/login` | Login, returns JWT token |

### Users *(requires Bearer token)*
| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/users` | List all users |
| `GET` | `/api/users/:id` | Get user by ID |
| `POST` | `/api/users` | Create user |
| `PUT` | `/api/users/:id` | Update user |
| `DELETE` | `/api/users/:id` | Delete user |

---

## 📦 Production Build

Build the Vue frontend and run everything from a single Go binary:

```bash
# 1. Build frontend → outputs to static/
cd web && npm run build

# 2. Start Go server (serves frontend + API on same port)
cd .. && go run main.go
```

| URL | Serves |
|---|---|
| `http://localhost:8080/` | Vue SPA (dashboard) |
| `http://localhost:8080/dashboard` | Vue Router handles it |
| `http://localhost:8080/api/...` | Go REST API |

---

## 🐳 Docker

```bash
docker-compose up -d
```

---

## 📝 License

MIT
