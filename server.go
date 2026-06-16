package main
// Trigger air reload for redesigned assets

import (
	"crypto/subtle"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

//go:embed dist/*
var embedDist embed.FS

// basicAuthMiddleware protects routes using SERVER_USERNAME and SERVER_PASSWORD
func basicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle CORS preflight
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		user, pass, ok := r.BasicAuth()

		expectedUser := os.Getenv("SERVER_USERNAME")
		if expectedUser == "" {
			expectedUser = "admin"
		}
		expectedPass := os.Getenv("SERVER_PASSWORD")
		if expectedPass == "" {
			expectedPass = "admin"
		}

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(expectedUser)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(expectedPass)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Admin Dashboard"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// StartAdminServer runs the dashboard API and serves frontend files
func StartAdminServer(port string) {
	mux := http.NewServeMux()

	// API Handlers (wrapped in basic auth)
	apiMux := http.NewServeMux()

	// Auth verification
	apiMux.HandleFunc("GET /api/auth/verify", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Routes CRUD
	apiMux.HandleFunc("GET /api/routes", handleGetRoutes)
	apiMux.HandleFunc("POST /api/routes", handleCreateRoute)
	apiMux.HandleFunc("PUT /api/routes/{id}", handleUpdateRoute)
	apiMux.HandleFunc("DELETE /api/routes/{id}", handleDeleteRoute)

	// Blocklist CRUD
	apiMux.HandleFunc("GET /api/blocklist", handleGetBlocklist)
	apiMux.HandleFunc("POST /api/blocklist", handleCreateBlocklist)
	apiMux.HandleFunc("DELETE /api/blocklist/{id}", handleDeleteBlocklist)

	// Logs
	apiMux.HandleFunc("GET /api/logs", handleGetLogs)
	apiMux.HandleFunc("DELETE /api/logs", handleClearLogs)

	// Metrics / Dashboard
	apiMux.HandleFunc("GET /api/metrics", handleGetMetrics)
	apiMux.HandleFunc("DELETE /api/metrics", handleClearMetrics)
	apiMux.HandleFunc("GET /api/metrics/stream", handleMetricsSSE)

	// Per-route traffic SSE stream
	apiMux.HandleFunc("GET /api/routes/{id}/stream", handleRouteStreamSSE)

	// Register API with Basic Auth
	mux.Handle("/api/", basicAuthMiddleware(apiMux))

	// Serve Static Files (with SPA fallback & Local Dev proxy support)
	if os.Getenv("GO_ENV") == "local" {
		targetURL, err := url.Parse("http://localhost:5173")
		if err != nil {
			log.Fatalf("invalid local Vite url: %v", err)
		}
		devProxy := httputil.NewSingleHostReverseProxy(targetURL)

		mux.Handle("/", basicAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If path has /api/, bypass (safety fallback)
			if strings.HasPrefix(r.URL.Path, "/api/") {
				return
			}
			devProxy.ServeHTTP(w, r)
		})))
	} else {
		subFS, err := fs.Sub(embedDist, "dist")
		if err != nil {
			log.Fatalf("failed to sub embedDist FS: %v", err)
		}
		fileServer := http.FileServer(http.FS(subFS))

		mux.Handle("/", basicAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If path has /api/, bypass (safety fallback)
			if strings.HasPrefix(r.URL.Path, "/api/") {
				return
			}

			// Clean and check if path exists in subFS
			path := strings.TrimPrefix(r.URL.Path, "/")
			if path == "" {
				path = "index.html"
			}

			_, err := subFS.Open(path)
			if err != nil {
				// File does not exist, serve index.html
				indexFile, err := subFS.Open("index.html")
				if err != nil {
					http.Error(w, "index.html not found in embedded assets", http.StatusNotFound)
					return
				}
				defer indexFile.Close()
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				_, _ = io.Copy(w, indexFile)
				return
			}

			fileServer.ServeHTTP(w, r)
		})))
	}

	server := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: mux,
	}

	log.Printf("[server] admin dashboard listening on port %s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("admin server failed: %v", err)
	}
}

// ─── Route Handlers ───────────────────────────────────────────────────────────

func handleGetRoutes(w http.ResponseWriter, r *http.Request) {
	var routes []ProxyRoute
	if err := db.Order("domain ASC").Find(&routes).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(routes)
}

func handleCreateRoute(w http.ResponseWriter, r *http.Request) {
	var route ProxyRoute
	if err := json.NewDecoder(r.Body).Decode(&route); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if route.Domain == "" {
		http.Error(w, "domain is required", http.StatusBadRequest)
		return
	}

	if err := db.Create(&route).Error; err != nil {
		http.Error(w, "failed to create route: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Reload dynamic config cache
	reloadConfigCache()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(route)
}

func handleUpdateRoute(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	var route ProxyRoute
	if err := db.First(&route, id).Error; err != nil {
		http.Error(w, "route not found", http.StatusNotFound)
		return
	}

	var updated ProxyRoute
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	route.Domain = updated.Domain
	route.SchemaType = updated.SchemaType
	route.TargetURL = updated.TargetURL
	route.DynamicResolveURL = updated.DynamicResolveURL
	route.UseBasicAuth = updated.UseBasicAuth
	route.BasicAuthUsername = updated.BasicAuthUsername
	route.BasicAuthPassword = updated.BasicAuthPassword
	route.UseValidationMiddleware = updated.UseValidationMiddleware
	route.ValidationMiddlewareURL = updated.ValidationMiddlewareURL
	route.SSLActive = updated.SSLActive
	route.LogPathPrefix = updated.LogPathPrefix
	route.RateLimit = updated.RateLimit
	route.RateLimitUnit = updated.RateLimitUnit
	route.RateLimitMethod = updated.RateLimitMethod
	route.RateLimitHeaderKey = updated.RateLimitHeaderKey
	route.RateLimitHeaderValue = updated.RateLimitHeaderValue

	if err := db.Save(&route).Error; err != nil {
		http.Error(w, "failed to update route: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Reload dynamic config cache
	reloadConfigCache()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(route)
}

func handleDeleteRoute(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	if err := db.Delete(&ProxyRoute{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Reload dynamic config cache
	reloadConfigCache()

	w.WriteHeader(http.StatusNoContent)
}

// ─── Blocklist Handlers ───────────────────────────────────────────────────────

func handleGetBlocklist(w http.ResponseWriter, r *http.Request) {
	var list []BlockedIP
	if err := db.Order("created_at DESC").Find(&list).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func handleCreateBlocklist(w http.ResponseWriter, r *http.Request) {
	var item BlockedIP
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if item.IPAddress == "" {
		http.Error(w, "IP address is required", http.StatusBadRequest)
		return
	}

	if err := db.Create(&item).Error; err != nil {
		http.Error(w, "failed to block IP: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func handleDeleteBlocklist(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	if err := db.Delete(&BlockedIP{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ─── Logs Handlers ────────────────────────────────────────────────────────────

func handleGetLogs(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	var logs []ProxyLog
	if err := db.Order("timestamp DESC").Limit(limit).Find(&logs).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func handleClearLogs(w http.ResponseWriter, r *http.Request) {
	if err := db.Where("1 = 1").Delete(&ProxyLog{}).Error; err != nil {
		http.Error(w, "failed to clear logs: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func handleClearMetrics(w http.ResponseWriter, r *http.Request) {
	if err := db.Where("1 = 1").Delete(&ProxyMetric{}).Error; err != nil {
		http.Error(w, "failed to clear metrics: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// ─── Metrics / Dashboard ──────────────────────────────────────────────────────

type DashboardMetrics struct {
	TotalRequests      int64                  `json:"total_requests"`
	SuccessRequests    int64                  `json:"success_requests"` // 2xx
	ErrorRequests      int64                  `json:"error_requests"`   // 4xx, 5xx
	AverageLatencyMs   float64                `json:"average_latency_ms"`
	VolumeSeries       []ChartDataPoint       `json:"volume_series"`
	LatencySeries      []ChartDataPoint       `json:"latency_series"`
	StatusCodesSeries  []PieChartPoint        `json:"status_codes_series"`
}

type ChartDataPoint struct {
	Timestamp int64   `json:"timestamp"` // Unix timestamp in Ms
	Value     float64 `json:"value"`
}

type PieChartPoint struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

func handleGetMetrics(w http.ResponseWriter, r *http.Request) {
	var totalRequests int64
	db.Model(&ProxyLog{}).Count(&totalRequests)

	var successRequests int64
	db.Model(&ProxyLog{}).Where("status_code >= 200 AND status_code < 300").Count(&successRequests)

	var errorRequests int64
	db.Model(&ProxyLog{}).Where("status_code >= 400 OR status_code == 0").Count(&errorRequests)

	var avgLatency float64
	db.Model(&ProxyLog{}).Select("COALESCE(AVG(response_time_ms), 0)").Row().Scan(&avgLatency)

	// Fetch logs for graphs (e.g., last 24 hours)
	var logs []ProxyLog
	now := time.Now()
	twentyFourHoursAgo := now.Add(-24 * time.Hour)
	db.Where("timestamp >= ?", twentyFourHoursAgo).Order("timestamp ASC").Find(&logs)

	// Status code distribution
	statusMap := make(map[string]int64)
	for _, l := range logs {
		codeStr := strconv.Itoa(l.StatusCode)
		if l.StatusCode == 0 {
			codeStr = "Failed/Blocked"
		}
		statusMap[codeStr]++
	}

	statusCodesSeries := make([]PieChartPoint, 0, len(statusMap))
	for k, v := range statusMap {
		statusCodesSeries = append(statusCodesSeries, PieChartPoint{Name: k, Value: v})
	}

	// Volume and latency graphs (aggregated hourly)
	hourlyData := make(map[time.Time][]ProxyLog)
	for _, l := range logs {
		hour := l.Timestamp.Truncate(time.Hour)
		hourlyData[hour] = append(hourlyData[hour], l)
	}

	volumeSeries := make([]ChartDataPoint, 0)
	latencySeries := make([]ChartDataPoint, 0)

	// Sort hours
	for h := twentyFourHoursAgo.Truncate(time.Hour); h.Before(now.Add(time.Hour)); h = h.Add(time.Hour) {
		hLogs := hourlyData[h]
		count := len(hLogs)
		var totalLat int64
		for _, l := range hLogs {
			totalLat += l.ResponseTimeMs
		}
		avgLat := float64(0)
		if count > 0 {
			avgLat = float64(totalLat) / float64(count)
		}

		ts := h.UnixNano() / int64(time.Millisecond)
		volumeSeries = append(volumeSeries, ChartDataPoint{Timestamp: ts, Value: float64(count)})
		latencySeries = append(latencySeries, ChartDataPoint{Timestamp: ts, Value: avgLat})
	}

	metrics := DashboardMetrics{
		TotalRequests:     totalRequests,
		SuccessRequests:   successRequests,
		ErrorRequests:     errorRequests,
		AverageLatencyMs:   avgLatency,
		VolumeSeries:       volumeSeries,
		LatencySeries:      latencySeries,
		StatusCodesSeries:  statusCodesSeries,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// handleMetricsSSE streams dashboard metrics via Server-Sent Events
func handleMetricsSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	ctx := r.Context()

	sendMetrics := func() bool {
		var totalRequests int64
		db.Model(&ProxyLog{}).Count(&totalRequests)

		var successRequests int64
		db.Model(&ProxyLog{}).Where("status_code >= 200 AND status_code < 300").Count(&successRequests)

		var errorRequests int64
		db.Model(&ProxyLog{}).Where("status_code >= 400 OR status_code == 0").Count(&errorRequests)

		var avgLatency float64
		db.Model(&ProxyLog{}).Select("COALESCE(AVG(response_time_ms), 0)").Row().Scan(&avgLatency)

		// Fetch ProxyMetric data for charts (last 24 hours)
		var metricsData []ProxyMetric
		cutoff := time.Now().Add(-24 * time.Hour)
		db.Where("show_at >= ?", cutoff).Order("show_at ASC").Find(&metricsData)

		volumeSeries := make([]ChartDataPoint, 0, len(metricsData))
		latencySeries := make([]ChartDataPoint, 0, len(metricsData))
		for _, m := range metricsData {
			ts := m.ShowAt.UnixNano() / int64(time.Millisecond)
			volumeSeries = append(volumeSeries, ChartDataPoint{Timestamp: ts, Value: float64(m.RequestVolume)})
			latencySeries = append(latencySeries, ChartDataPoint{Timestamp: ts, Value: m.RequestLatency})
		}

		metrics := map[string]interface{}{
			"total_requests":     totalRequests,
			"success_requests":   successRequests,
			"error_requests":     errorRequests,
			"average_latency_ms": avgLatency,
			"volume_series":      volumeSeries,
			"latency_series":     latencySeries,
		}

		data, err := json.Marshal(metrics)
		if err != nil {
			return false
		}

		_, err = fmt.Fprintf(w, "data: %s\n\n", data)
		if err != nil {
			return false
		}
		flusher.Flush()
		return true
	}

	if !sendMetrics() {
		return
	}

	for {
		select {
		case <-ticker.C:
			if !sendMetrics() {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

// handleRouteStreamSSE streams live request logs for a specific route via broker
func handleRouteStreamSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	var route ProxyRoute
	if err := db.First(&route, id).Error; err != nil {
		http.Error(w, "route not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ctx := r.Context()
	ch := make(chan []byte, 64)
	sseBroker.register <- subscription{ch: ch, domain: route.Domain}

	defer func() {
		sseBroker.unregister <- ch
	}()

	// Send initial comment to confirm connection
	fmt.Fprintf(w, ": connected to %s\n\n", route.Domain)
	flusher.Flush()

	for {
		select {
		case data := <-ch:
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		case <-ctx.Done():
			return
		}
	}
}
