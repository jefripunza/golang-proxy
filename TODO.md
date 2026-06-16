# Task List

## Backlog
- Overview gunakan SSE (Server Send Event) untuk realtime
- Proxy Routes > di antara list, jika di klik akan muncul dialog modal untuk xterm / terminal web yang realtime membaca network yang jalan di proxy itu (SSE nya akan konek ketika dialog nya muncul, jika tidak maka SSE nya di disconnect)
- dark mode nya belum work, revisi sampai work ya

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
