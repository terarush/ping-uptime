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

## 🚀 Getting Started (Docker Installation)

The easiest way to run `ping-uptime` is using Docker.

### 1. Create a Configuration File

Create a `.env` file in your working directory and adjust the configurations as needed.

```env
APP_NAME=ping-uptime
SERVER_MODE=info
PORT=8080
HTTP_TIMEOUT=60

# Database Configuration
DB_DRIVER=mysql         # mysql | postgres
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=ping_uptime
DB_USERNAME=root
DB_PASSWORD=your_password

# Connection Pool (Optional)
POOL_CONN_IDLE=200
POOL_CONN_MAX=300
POOL_CONN_LIFETIME=60

# JWT Authentication
JWT_SIGNATURE_KEY=your_secret_key
JWT_DAY_EXPIRED=60

# Discord Notifications (Webhook)
DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/your/webhook
```

### 2. Run with Docker

Run the container using the following command, passing your `.env` file:

```bash
docker run -d \
  --name ping-uptime \
  -p 8080:8080 \
  --env-file .env \
  rafia9005/ping-uptime:latest
```

Once running, you can access the application:
- **Dashboard**: `http://localhost:8080`
- **REST API**: `http://localhost:8080/api`

---

## 🐳 Docker Compose

Alternatively, you can run `ping-uptime` along with a MySQL database using `docker-compose`.

1. Make sure you have `.env` configured.
2. Run:
```bash
docker-compose up -d
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

## 📝 License

MIT
