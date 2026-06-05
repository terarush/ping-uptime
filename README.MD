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

## 🚀 Getting Started (Docker Installation)

The easiest way to run `ping-uptime` is using Docker.

## Run on Linux System

Run the ping-uptime on your Linux system without Docker.

```bash
curl -fsSL https://raw.githubusercontent.com/terarush/ping-uptime/main/install.sh | bash
```

### Run with Docker

Run the container using the following command, passing your `.env` file:

```bash
docker run -d \
  --name ping-uptime \
  -p 8080:8080 \
  rafia9005/ping-uptime:latest
```

Run the container with custom volume mounts for configuration:

```bash
docker run -d \
  --name ping-uptime \
  -p 8081:8080 \
  -e DB_NAME=/app/data/ping-uptime \
  -v $(pwd)/data:/app/data \
  rafia9005/ping-uptime:latest
```
