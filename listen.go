package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// getClientIP extracts the host IP from RemoteAddr
func getClientIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// isIPBlocked checks if the client IP is present in the blocklist
func isIPBlocked(ip string) bool {
	var count int64
	db.Model(&BlockedIP{}).Where("ip_address = ?", ip).Count(&count)
	return count > 0
}

// resolveRouteDB matches the incoming Host+Path against routeConfigDB
func resolveRouteDB(r *http.Request) (ProxyRoute, bool) {
	routeMu.RLock()
	defer routeMu.RUnlock()

	host := r.Host
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}

	path := strings.TrimPrefix(r.URL.Path, "/")
	full := host + "/" + path
	reqPath := "/" + path

	for _, key := range routeKeys {
		// Normalize key: strip port for host comparison
		keyHost := key
		keyFull := key
		if h, p, err := net.SplitHostPort(key); err == nil {
			_ = p
			keyHost = h
			if idx := strings.Index(key, "/"); idx > 0 {
				keyFull = keyHost + key[idx:]
			} else {
				keyFull = keyHost
			}
		}

		// Full host+path prefix match (using normalized key)
		if strings.HasPrefix(full, keyFull) {
			return routeConfigDB[key], true
		}

		// Exact host match (port-stripped)
		if keyHost == host {
			return routeConfigDB[key], true
		}

		// Path-only match: extract path from key, match incoming path regardless of host
		if idx := strings.Index(key, "/"); idx > 0 {
			keyPath := key[idx:]
			if strings.HasPrefix(reqPath, keyPath+"/") || reqPath == keyPath {
				return routeConfigDB[key], true
			}
		}
	}
	return ProxyRoute{}, false
}

// dynamicRouteResponse represents the JSON payload from a dynamic resolve URL
type dynamicRouteResponse struct {
	URL       string     `json:"url"`
	BasicAuth *BasicAuth `json:"basicAuth,omitempty"`
}

// resolveDynamicRoute queries the dynamic resolver endpoint
func resolveDynamicRoute(resolveURL string) (string, *BasicAuth, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(resolveURL)
	if err != nil {
		return "", nil, fmt.Errorf("dynamic resolver request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("dynamic resolver returned status %d", resp.StatusCode)
	}

	var res dynamicRouteResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", nil, fmt.Errorf("failed to decode dynamic resolver JSON: %w", err)
	}

	if res.URL == "" {
		return "", nil, fmt.Errorf("dynamic resolver returned empty target URL")
	}

	return res.URL, res.BasicAuth, nil
}

// validateRequestMiddleware forwards request context to the validation URL
func validateRequestMiddleware(validationURL string, r *http.Request) (bool, error) {
	// Read and buffer the body so it can be forwarded and then re-read later by the proxy
	var bodyBytes []byte
	if r.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(r.Body)
		if err != nil {
			return false, err
		}
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	client := &http.Client{Timeout: 5 * time.Second}
	valReq, err := http.NewRequest(r.Method, validationURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return false, err
	}

	// Copy headers
	for k, vv := range r.Header {
		for _, v := range vv {
			valReq.Header.Add(k, v)
		}
	}

	valResp, err := client.Do(valReq)
	if err != nil {
		return false, err
	}
	defer valResp.Body.Close()

	if valResp.StatusCode >= 200 && valResp.StatusCode < 300 {
		return true, nil
	}

	return false, fmt.Errorf("validation returned status %d", valResp.StatusCode)
}

// logProxyRequest saves a proxy request log asynchronously
func logProxyRequest(domain, path, method string, statusCode int, duration time.Duration, clientIP, errMsg string) {
	go func() {
		logEntry := ProxyLog{
			Timestamp:      time.Now(),
			Domain:         domain,
			Path:           path,
			Method:         method,
			StatusCode:     statusCode,
			ResponseTimeMs: duration.Milliseconds(),
			SourceIP:       clientIP,
			ErrorMessage:   errMsg,
		}
		if err := db.Create(&logEntry).Error; err != nil {
			log.Printf("[logger] failed to write proxy log: %v", err)
		}
	}()
}

// checkAuthProxy verifies authentication using Basic Auth
func checkAuthProxy(r *http.Request, auth *BasicAuth) bool {
	// 1. Session cookie
	for _, c := range r.Cookies() {
		if c.Name == "app_session" && c.Value == sessionToken {
			return true
		}
	}

	// 2. Basic Auth check
	if auth == nil || (auth.Username == "" && auth.Password == "") {
		return true
	}

	headerVal := r.Header.Get("Authorization")
	if strings.HasPrefix(headerVal, "Basic ") {
		decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(headerVal, "Basic "))
		if err == nil {
			parts := strings.SplitN(string(decoded), ":", 2)
			if len(parts) == 2 && parts[0] == auth.Username && parts[1] == auth.Password {
				return true
			}
		}
	}
	return false
}

// ─── HTTP proxy handler ───────────────────────────────────────────────────────

func handleProxy(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	clientIP := getClientIP(r)

	// 1. IP Blocklist validation
	if isIPBlocked(clientIP) {
		http.Error(w, "Forbidden: IP address blocked", http.StatusForbidden)
		logProxyRequest(r.Host, r.URL.Path, r.Method, http.StatusForbidden, time.Since(startTime), clientIP, "IP Blocked")
		return
	}

	// 2. Resolve Route configuration
	route, ok := resolveRouteDB(r)
	if !ok {
		msg := "no route matched for " + r.Host + r.URL.Path
		http.Error(w, msg, http.StatusBadGateway)
		logProxyRequest(r.Host, r.URL.Path, r.Method, http.StatusBadGateway, time.Since(startTime), clientIP, msg)
		return
	}

	targetURL := route.TargetURL
	var finalAuth *BasicAuth
	if route.UseBasicAuth {
		finalAuth = &BasicAuth{
			Username: route.BasicAuthUsername,
			Password: route.BasicAuthPassword,
		}
	}

	// 3. Dynamic route resolution
	if route.SchemaType == "dynamic" {
		var err error
		targetURL, finalAuth, err = resolveDynamicRoute(route.DynamicResolveURL)
		if err != nil {
			msg := "dynamic routing resolve failed: " + err.Error()
			http.Error(w, msg, http.StatusBadGateway)
			logProxyRequest(r.Host, r.URL.Path, r.Method, http.StatusBadGateway, time.Since(startTime), clientIP, msg)
			return
		}
	}

	// 4. Basic Auth validation
	if !checkAuthProxy(r, finalAuth) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Proxy"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		logProxyRequest(r.Host, r.URL.Path, r.Method, http.StatusUnauthorized, time.Since(startTime), clientIP, "Unauthorized")
		return
	}

	// Set session cookie on successful basic auth
	if r.Header.Get("Authorization") != "" {
		setSessionCookie(w)
	}

	// 5. Validation middleware
	if route.UseValidationMiddleware {
		ok, err := validateRequestMiddleware(route.ValidationMiddlewareURL, r)
		if !ok || err != nil {
			errMsg := "validation middleware failed"
			if err != nil {
				errMsg += ": " + err.Error()
			}
			http.Error(w, errMsg, http.StatusForbidden)
			logProxyRequest(r.Host, r.URL.Path, r.Method, http.StatusForbidden, time.Since(startTime), clientIP, errMsg)
			return
		}
	}

	// Clean up URL: targetURL might be "http://something.com" or just "something.com"
	cleanedTarget := targetURL
	if !strings.HasPrefix(cleanedTarget, "http://") && !strings.HasPrefix(cleanedTarget, "https://") {
		cleanedTarget = "http://" + cleanedTarget
	}

	// Parse host from URL
	parsedURL := strings.TrimPrefix(cleanedTarget, "http://")
	parsedURL = strings.TrimPrefix(parsedURL, "https://")
	if idx := strings.Index(parsedURL, "/"); idx != -1 {
		parsedURL = parsedURL[:idx]
	}

	outReq, err := http.NewRequest(r.Method, cleanedTarget+r.URL.RequestURI(), r.Body)
	if err != nil {
		msg := "bad gateway: " + err.Error()
		http.Error(w, msg, http.StatusBadGateway)
		logProxyRequest(r.Host, r.URL.Path, r.Method, http.StatusBadGateway, time.Since(startTime), clientIP, msg)
		return
	}
	outReq.Header = stripProxyHeaders(r.Header, parsedURL)

	resp, err := http.DefaultTransport.RoundTrip(outReq)
	if err != nil {
		msg := "bad gateway: " + err.Error()
		http.Error(w, msg, http.StatusBadGateway)
		logProxyRequest(r.Host, r.URL.Path, r.Method, http.StatusBadGateway, time.Since(startTime), clientIP, msg)
		return
	}
	defer resp.Body.Close()

	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)

	logProxyRequest(r.Host, r.URL.Path, r.Method, resp.StatusCode, time.Since(startTime), clientIP, "")
}

// ─── WebSocket / HTTP upgrade handler ────────────────────────────────────────

func upgradeHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	clientIP := getClientIP(r)

	// 1. IP Blocklist validation
	if isIPBlocked(clientIP) {
		if hj, ok := w.(http.Hijacker); ok {
			conn, _, _ := hj.Hijack()
			conn.Write([]byte("HTTP/1.1 403 Forbidden\r\n\r\n"))
			conn.Close()
		}
		logProxyRequest(r.Host, r.URL.Path, r.Method, http.StatusForbidden, time.Since(startTime), clientIP, "WS IP Blocked")
		return
	}

	// 2. Resolve Route configuration
	route, ok := resolveRouteDB(r)
	if !ok {
		if hj, ok := w.(http.Hijacker); ok {
			conn, _, _ := hj.Hijack()
			conn.Write([]byte("HTTP/1.1 502 Bad Gateway\r\n\r\n"))
			conn.Close()
		}
		logProxyRequest(r.Host, r.URL.Path, r.Method, http.StatusBadGateway, time.Since(startTime), clientIP, "WS Route not found")
		return
	}

	targetURL := route.TargetURL
	var finalAuth *BasicAuth
	if route.UseBasicAuth {
		finalAuth = &BasicAuth{
			Username: route.BasicAuthUsername,
			Password: route.BasicAuthPassword,
		}
	}

	// 3. Dynamic route resolution
	if route.SchemaType == "dynamic" {
		var err error
		targetURL, finalAuth, err = resolveDynamicRoute(route.DynamicResolveURL)
		if err != nil {
			if hj, ok := w.(http.Hijacker); ok {
				conn, _, _ := hj.Hijack()
				conn.Write([]byte("HTTP/1.1 502 Bad Gateway\r\n\r\n"))
				conn.Close()
			}
			logProxyRequest(r.Host, r.URL.Path, r.Method, http.StatusBadGateway, time.Since(startTime), clientIP, "WS Dynamic resolve failed: "+err.Error())
			return
		}
	}

	// 4. Basic Auth validation
	if !checkAuthProxy(r, finalAuth) {
		fmt.Fprintf(w.(io.Writer), "HTTP/1.1 401 Unauthorized\r\nWWW-Authenticate: Basic realm=\"Proxy\"\r\n\r\n")
		if hj, ok := w.(http.Hijacker); ok {
			conn, _, _ := hj.Hijack()
			conn.Close()
		}
		logProxyRequest(r.Host, r.URL.Path, r.Method, http.StatusUnauthorized, time.Since(startTime), clientIP, "WS Unauthorized")
		return
	}

	// 5. Validation middleware
	if route.UseValidationMiddleware {
		ok, err := validateRequestMiddleware(route.ValidationMiddlewareURL, r)
		if !ok || err != nil {
			if hj, ok := w.(http.Hijacker); ok {
				conn, _, _ := hj.Hijack()
				conn.Write([]byte("HTTP/1.1 403 Forbidden\r\n\r\n"))
				conn.Close()
			}
			logProxyRequest(r.Host, r.URL.Path, r.Method, http.StatusForbidden, time.Since(startTime), clientIP, "WS Validation failed")
			return
		}
	}

	// Clean up Target
	cleanedTarget := targetURL
	cleanedTarget = strings.TrimPrefix(cleanedTarget, "http://")
	cleanedTarget = strings.TrimPrefix(cleanedTarget, "https://")
	if idx := strings.Index(cleanedTarget, "/"); idx != -1 {
		cleanedTarget = cleanedTarget[:idx]
	}

	target, err := net.Dial("tcp", cleanedTarget)
	if err != nil {
		if hj, ok := w.(http.Hijacker); ok {
			conn, _, _ := hj.Hijack()
			conn.Write([]byte("HTTP/1.1 502 Bad Gateway\r\n\r\n"))
			conn.Close()
		}
		logProxyRequest(r.Host, r.URL.Path, r.Method, http.StatusBadGateway, time.Since(startTime), clientIP, "WS TCP Dial failed: "+err.Error())
		return
	}

	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "hijacking not supported", http.StatusInternalServerError)
		target.Close()
		logProxyRequest(r.Host, r.URL.Path, r.Method, http.StatusInternalServerError, time.Since(startTime), clientIP, "WS Hijack not supported")
		return
	}
	clientConn, _, err := hj.Hijack()
	if err != nil {
		target.Close()
		return
	}

	// Re-write the upgrade request to the target
	stripped := stripProxyHeaders(r.Header, cleanedTarget)
	var sb strings.Builder
	sb.WriteString(r.Method + " " + r.URL.RequestURI() + " HTTP/1.1\r\n")
	for k, vv := range stripped {
		for _, v := range vv {
			sb.WriteString(k + ": " + v + "\r\n")
		}
	}
	sb.WriteString("\r\n")
	target.Write([]byte(sb.String()))

	// Bidirectional pipe
	done := make(chan struct{}, 2)
	go func() { io.Copy(target, clientConn); done <- struct{}{} }()
	go func() { io.Copy(clientConn, target); done <- struct{}{} }()
	<-done

	target.Close()
	clientConn.Close()

	logProxyRequest(r.Host, r.URL.Path, r.Method, 101, time.Since(startTime), clientIP, "")
}

func setSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "app_session",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

var proxyHeadersToStrip = []string{
	"X-Forwarded-For",
	"X-Forwarded-Proto",
	"X-Forwarded-Host",
	"X-Forwarded-Port",
	"X-Forwarded-Scheme",
	"X-Forwarded-Ssl",
	"X-Real-Ip",
	"Forwarded",
	"Via",
	"Cf-Connecting-Ip",
	"True-Client-Ip",
	"Authorization",
}

func stripProxyHeaders(h http.Header, targetHost string) http.Header {
	out := h.Clone()
	for _, k := range proxyHeadersToStrip {
		out.Del(k)
	}
	out.Set("Host", targetHost)
	out.Set("Origin", "http://"+targetHost)
	if out.Get("Referer") != "" {
		out.Set("Referer", "http://"+targetHost+"/")
	}
	return out
}

