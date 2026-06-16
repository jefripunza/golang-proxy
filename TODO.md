# Task List

## Backlog
- di halaman Logs, button "Clear Logs", jika di execute maka yang clear hanya table log saja "ProxyLog" dan jangan clear table lain;
- di halaman Overview, tambahkan button clear untuk clear table "ProxyMetric";
- untuk rate limiter, bikin table baru saja, table rate limiter tidak boleh di delete / clear;
- 

## Error / Bug
- 

## Finish
- [x] Fix vue-doctor warnings and errors to achieve 100% Perfect audit score
- [x] Fix sidebar visibility width layout issue via inline styles
- [x] Fix sidebar squished layout issue by restoring default Tailwind spacing scale in global.css
- [x] Fully protect static files and assets with Basic Auth to prevent background dashboard rendering before login
- [x] Implement embedded static assets in Go binary via embed.FS
- [x] Implement local development mode asset reverse proxy to Vite server
- [x] Remove Vue custom Login view and integrate native browser basic auth flow
- [x] Fix runtime invalid memory address or nil pointer dereference panic during startup
- [x] Set up SQLite3 connection using GORM
- [x] Create GORM models for ProxyRoute, ProxyLog, BlockedIP in `models.go`
- [x] Implement dual concurrent HTTP server runner (LISTEN proxy, SERVER dashboard)
- [x] Implement Basic Authentication middleware for the SERVER using environment variables
- [x] Implement IP Blocklist validation (returns 403 Forbidden if blocked)
- [x] Implement Static Routing schema (source domain -> target URL)
- [x] Implement Dynamic Routing schema (resolves target by calling external backend resolver)
- [x] Implement Validation Middleware (forward request context, proceed on 2xx status)
- [x] Implement Auto SSL Let's Encrypt manager via autocert
- [x] Configure Tailwind CSS dark mode settings based on DESIGN.md
- [x] Create header component with a functioning Dark Mode Toggle
- [x] Create reusable Modal/Dialog component (outside click ignored)
- [x] Create reusable Chart component using highcharts
- [x] Create reusable Form/FormField component wrapper and Input fields with live onchange Zod validation
- [x] Set up Axios service client with basic auth interceptor
- [x] Implement navigation layout (Overview, Proxy Routes, IP Block List, Activity Logs, Settings)
- [x] Build visual analytics dashboard using Highcharts
- [x] Fix dark mode toggle — localStorage persistence + light theme CSS
- [x] Overview SSE real-time — replace polling with Server-Sent Events for live stats
- [x] Proxy Routes terminal — click route to open xterm.js live traffic monitor via per-route SSE
- [x] Restore left sidebar layout (remove horizontal pill nav)
- [x] Remove footer trust bar
- [x] Remove "Product" nav placeholder
- [x] Rename title to "Golang Proxy"
- [x] Remove Settings page and route
- [x] Fix Overview SSE with EventSource + withCredentials
- [x] Polish announcement bar — double-oval icon + "THE BEST NEVER GUESS" text
- [x] Move logo (abstract right-arrow) to sidebar top, simplify header to just theme toggle + avatar
- [x] Sidebar nav: active pill-outline style (blue border + subtle bg, rounded-full)
- [x] Simplify nav labels: Overview / Routes / Block List / Logs
- [x] Update page headings to match nav labels
- [x] Sidebar: burger button to collapse/expand sidebar
- [x] Header: remove Admin text + avatar, keep only burger + dark mode toggle
- [x] Overview: remove hero card (text + refresh/view buttons + preview mockup)
- [x] Routes: terminal button in Actions column with console icon
- [x] Xterm: auto-resize via ResizeObserver, fills 80vw dialog
- [x] Routes: Log URL Prefix optional field — filter proxy request logging by path prefix
- [x] Log URL Prefix: support comma-separated multiple prefixes with trim spaces
- [x] Logs: Clear Logs button + DELETE /api/logs endpoint
- [x] Logs detail: add request body + response body (capped 64KB, labeled "[file/large body — truncated]" if oversized)
- [x] Dockerfile + docker.sh adapted for golang-proxy (bun FE + go BE, multi-arch buildx)
- [x] Rate Limiter: configurable per-route (limit, unit second/minute/hour/day, method compact/ip/header)
- [x] Compact rate limiter: serialize all request headers to JSON, hash with UUIDv5 for unique fingerprint
- [x] Fix route port matching: route with port (e.g. localhost:8080) must match exact port, not just host
- [x] ProxyMetric table: per-minute request_volume + request_latency, worker inserts at HH:MM:00
- [x] Overview charts: SSE streams ProxyMetric data for real-time Volume + Latency charts
