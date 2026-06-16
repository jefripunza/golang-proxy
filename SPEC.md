# Project Specification (SPEC.md)

This document outlines the technical stack, design requirements, and feature scope for the **Dynamic Golang Proxy** project.

---

## 💻 Tech Stack

### Frontend (FE)
- **Framework**: Vue 3 + TypeScript (Composition API, `<script setup>` pattern preferred).
- **Styling**: Tailwind CSS (v4.3.1).
- **State Management**: Pinia.
- **Data Fetching**: TanStack Query (Vue Query), using endpoint functions imported from the services layer.
- **HTTP Client**: Axios (located in the `services/` directory, e.g., `src/services/`).

### Backend (BE)
- **Runtime**: Go (Golang).
- **Web Server**: Standard library `net/http` package (plain/standard HTTP, no custom heavy framework).
- **Architecture**: Runs **two separate HTTP servers** concurrently:
  - **LISTEN**: The HTTP server/service handling the dynamic reverse proxy traffic, routing requests, and managing Auto SSL Let's Encrypt certificates.
  - **SERVER**: The HTTP server/service handling the web dashboard API requests and administration.
    - **Security**: The entire SERVER (both static asset handlers and administrative API endpoints) must be fully protected by **Basic Authentication** loaded from environment variables (`SERVER_USERNAME` and `SERVER_PASSWORD`). Default credentials are `admin` / `admin`. If the credentials are not provided or invalid, no files, HTML structures, or assets should be returned to the client to prevent the page structure from being visible prior to authentication.
    - **Frontend Assets**:
      - Served from embedded assets compiled inside the Go binary:
        ```go
        //go:embed dist/*
        var embedDist embed.FS
        ```
      - Fallback requests (instead of returning 404) must serve `index.html` from this embedded filesystem to support SPA routing.
      - **Development Mode**: If the environment variable `GO_ENV=local` is set, the server must bypass `embedDist` and instead reverse proxy all frontend asset requests to the local Vite development server running at `http://localhost:5173`.
- **ORM**: GORM.
- **Database**: SQLite3.

---

## 🎨 UI/UX & Design Guidelines

### High Quality UI & Design System
- **No "AI Slop"**: All UI elements must be premium, professional, responsive, and follow modern design principles.
- **Design System Consistency**: Always read `DESIGN.md` in the root of the workspace before creating or modifying UI to ensure consistent colors, typography, spacing, and styling.
- **CRITICAL RULE**: If `DESIGN.md` **does not exist**, you **MUST STOP IMMEDIATELY** and request the user to visit [styles.refero.design](https://styles.refero.design/) to select a design context and create/provide `DESIGN.md` before proceeding.
- **UX Tooling**: Utilize the `.agents/skills/ui-ux-pro-max` skill to design elegant layouts, color palettes, and interactive components.
- **Dark Mode**: The application must support Dark Mode, with the toggle switch located clearly in the dashboard header.
- **Cross-Browser Compatibility**: Ensure Tailwind CSS styling and layouts are fully tested and run consistently across all modern web browsers.

### Frontend Component Architecture & Requirements
- **Highly Reusable Components**: All frontend code must be broken down into small, single-responsibility, reusable components. Do not build massive, monolithic page components.
- **Professional Folder Structure**: Components must be structured cleanly in dedicated directories (e.g., `components/common/`, `components/dashboard/`, `components/icons/`) inside the `src/` directory.
- **Modal Dialogs**: Dialogs and modal popups **must not close** when clicking outside the modal content area (outside clicks must be ignored).
- **Charts**: Always use `highcharts` and `highcharts-vue` for rendering charts. All chart instances must be wrapped inside reusable components.
- **Icons**: Always use `vue-icons-plus` as the standard icon library.
- **Forms & Zod Validation**:
  - All forms must use **Zod** for schema validation.
  - Validation must be double-checked (both at the form level and individual field level).
  - Create/use a reusable wrapper `<Form>` component to handle Zod validation.
  - Every form input field must be its own component with live validation (triggering and showing validation feedback immediately `onchange` or on input).

---

## 🎯 Project Scope & Key Features
This project implements a dynamic reverse proxy server:
1. **Proxy Routing Schemas**:
   - **Static**: Directly forwards traffic from a defined source to a defined target URL.
   - **Dynamic**: Resolves the target destination by calling an external backend. The backend returns a JSON response:
     - `url`: The target URL to proxy the request to.
     - `basicAuth` *(Optional)*: An object containing `username` and `password` to authenticate against the target.
2. **Validation Middleware**:
   - Both static and dynamic routing configurations can optionally define a validation middleware endpoint.
   - Before forwarding to the target, the proxy sends the entire request context (headers, body, etc.) to this validation endpoint.
   - If the validation endpoint returns a `2xx` status code, the proxy proceeds to the target.
   - If it returns any other status, execution halts at the middleware and the error is returned.
3. **IP Blocklist**:
   - The proxy engine (`LISTEN`) must validate incoming requests against a database-driven IP Blocklist.
   - If the request sender's source IP is in the blocklist, the request is immediately rejected with a `403 Forbidden` response and processing halts.
4. **Database & Storage**: All proxy routes, configurations, blocked IPs, and logs are saved in a local SQLite3 database.
5. **Dashboard Setup**: A web dashboard allows users to manage routes, configure the IP blocklist, and monitor proxy traffic metrics.
6. **Auto SSL (Let's Encrypt)**: Automatically manages Let's Encrypt SSL certificates for configured domains.
