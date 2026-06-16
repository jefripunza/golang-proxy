<p align="center">
  <img src="https://img.shields.io/badge/go-1.25-00ADD8?style=flat&logo=go" alt="Go">
  <img src="https://img.shields.io/badge/vue-3.x-4FC08D?style=flat&logo=vue.js" alt="Vue">
  <img src="https://img.shields.io/badge/sqlite-3-003B57?style=flat&logo=sqlite" alt="SQLite">
  <img src="https://img.shields.io/badge/docker-ready-2496ED?style=flat&logo=docker" alt="Docker">
</p>

# Golang Proxy

A high-performance **dynamic reverse proxy** written in Go with a professional Vue 3 web dashboard. Route traffic with static or dynamic targets, enforce IP blocklists, validate requests through middleware, and manage Let's Encrypt SSL certificates — all from a single binary.

---

## Features

- **Dual HTTP Servers** — LISTEN (proxy traffic) and SERVER (admin dashboard) run concurrently on separate ports.
- **Static Routing** — Forward traffic from a domain/path prefix directly to a target URL.
- **Dynamic Routing** — Resolve target destinations at runtime by calling an external backend resolver.
- **Validation Middleware** — Optionally forward request context to a validation endpoint; proceed only on `2xx`.
- **IP Blocklist** — Database-driven; blocked IPs immediately receive `403 Forbidden`.
- **Auto SSL** — Let's Encrypt integration via `autocert` with per-domain toggle.
- **Web Dashboard** — Manage routes, blocklist, view activity logs, and monitor metrics with real-time charts.
- **Live Traffic Terminal** — Click any route to open an xterm.js terminal streaming proxy traffic via SSE.
- **Request Inspection** — Click any log entry to see full request/response headers and body.
- **Dark / Light Mode** — Toggle with localStorage persistence.
- **Single Binary** — Frontend assets embedded via `embed.FS`. No external web server needed.
- **Multi-arch Docker** — `linux/amd64` + `linux/arm64` via buildx.

---

## Quick Start

### Docker

```bash
docker run -d \
  -p 80:80 \
  -p 443:443 \
  -p 8080:8080 \
  -p 8000:8000 \
  -e SERVER_USERNAME=admin \
  -e SERVER_PASSWORD=admin \
  --name golang-proxy \
  jefriherditriyanto/golang-proxy:latest
```

Dashboard: `http://localhost:8000` (login: `admin` / `admin`)
Proxy listener: `http://localhost:8080`

### From Source

```bash
# Frontend
bun install
bun run build

# Backend
go build -o go-proxy .
./go-proxy
```

Requires **Go 1.22+** and **Bun** (or Node 20+).

---

## Environment Variables

| Variable          | Default | Description                                                                    |
| ----------------- | ------- | ------------------------------------------------------------------------------ |
| `LISTEN_PORT`     | `8080`  | Port for incoming proxy traffic                                                |
| `SERVER_PORT`     | `8000`  | Port for admin dashboard + API                                                 |
| `SERVER_USERNAME` | `admin` | Basic auth username for dashboard                                              |
| `SERVER_PASSWORD` | `admin` | Basic auth password for dashboard |
| `GO_ENV` | — | Set to `local` to proxy frontend assets to Vite dev server at `localhost:5173` |

> SSL is enabled **per-route** via the `ssl_active` field in the Routes dashboard. The Let's Encrypt infrastructure (ports 80/443) runs automatically — just toggle SSL on any route and certs are provisioned.

---

## Architecture

```
                     ┌─────────────────┐
                     │   Golang Proxy  │
                     │                 │
  Client ──────────► │ LISTEN :8080    │ ──► Target Server
                     │  • Route match  │
                     │  • IP blocklist │
                     │  • Validation   │
                     │  • SSL (ACME)   │
                     │                 │
  Admin  ──────────► │ SERVER :8000    │
                     │  • Dashboard UI │
                     │  • REST API     │
                     │  • SSE Streams  │
                     │  • Basic Auth   │
                     │                 │
                     │  SQLite (GORM)  │
                     └─────────────────┘
```

---

## Tech Stack

| Layer            | Technology                          |
| ---------------- | ----------------------------------- |
| **Proxy Engine** | Go 1.24, `net/http`, `autocert`     |
| **ORM**          | GORM + SQLite3                      |
| **Frontend**     | Vue 3, TypeScript, Vite             |
| **Styling**      | Tailwind CSS v4                     |
| **State**        | Pinia                               |
| **Charts**       | Highcharts                          |
| **Icons**        | vue-icons-plus                      |
| **Terminal**     | xterm.js                            |
| **Validation**   | Zod                                 |
| **Dev Tools**    | ESLint, Oxlint, vue-tsc, vue-doctor |

---

## Dashboard Pages

| Page           | Description                                                                                |
| -------------- | ------------------------------------------------------------------------------------------ |
| **Overview**   | Real-time metrics (SSE): total requests, success rate, latency, volume/status charts       |
| **Routes**     | CRUD proxy routes (static/dynamic), Log URL Prefix filter, live traffic terminal per route |
| **Block List** | Add/remove IPs from the blocklist with reason tracking                                     |
| **Logs**       | Request log stream with detail inspector (headers, body, status, latency) and clear-all    |

---

## License

MIT
