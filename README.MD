# ping-uptime

Self-hosted uptime monitoring service built with Go (Echo v4) and Vue 3.
Monitor websites, APIs, and services — get Discord notifications on status changes.

---

## Features

- HTTP/HTTPS uptime monitoring — ping endpoints at configurable intervals
- Discord webhook notifications — alerts on down/recovered status changes
- Dashboard — real-time status, response time, uptime history
- JWT authentication — secure API and dashboard access
- Multi-database support — MySQL or PostgreSQL via GORM
- Modular architecture — domain-driven design, easy to extend
- Single binary — Go backend serves Vue frontend from static/

---

## Getting Started (Docker)

Easiest way to run ping-uptime is with Docker.

### Run on Linux (no Docker)

```bash
curl -fsSL https://raw.githubusercontent.com/terarush/ping-uptime/main/install.sh | bash
```

### Run with Docker

```bash
docker run -d \
  --name ping-uptime \
  -p 8080:8080 \
  rafia9005/ping-uptime:latest
```

With custom volume mounts:

```bash
docker run -d \
  --name ping-uptime \
  -p 8081:8080 \
  -e DB_NAME=/app/data/ping-uptime \
  -v $(pwd)/data:/app/data \
  rafia9005/ping-uptime:latest
```
