package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/acme/autocert"
)

func hostPolicy(ctx context.Context, host string) error {
	var count int64
	// Search for exact match or matches in the DB with SSL active
	err := db.Model(&ProxyRoute{}).Where("domain = ? AND ssl_active = ?", host, true).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return fmt.Errorf("host %q not configured for SSL", host)
}

func main() {
	// Initialize SQLite DB connection & migration
	initDB()

	// Start SSE log broker
	go sseBroker.run()

	// Start rate limiter cleanup
	go rl.cleanup()

	// Start metrics worker
	startMetricWorker()

	// Start auto-clear worker
	startAutoClearWorker()

	// Load route config cache from database
	reloadConfigCache()

	// Start concurrent Admin Dashboard Server (SERVER)
	adminPort := os.Getenv("SERVER_PORT")
	if adminPort == "" {
		adminPort = "8000"
	}
	go StartAdminServer(adminPort)

	// Proxy Handler with WebSocket / HTTP Upgrade detection
	proxyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// WebSocket / HTTP upgrade
		if strings.EqualFold(r.Header.Get("Upgrade"), "websocket") ||
			strings.EqualFold(r.Header.Get("Connection"), "upgrade") {
			upgradeHandler(w, r)
			return
		}
		handleProxy(w, r)
	})

	// Let's Encrypt Auto SSL manager (per-route via ssl_active field)
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: hostPolicy,
		Cache:      autocert.DirCache("certs"),
	}

	// Port 80 handles HTTP-01 ACME challenges
	go func() {
		log.Println("[proxy] HTTP-01 challenge listener on port 80")
		if err := http.ListenAndServe("0.0.0.0:80", certManager.HTTPHandler(nil)); err != nil {
			log.Fatalf("HTTP challenge listener failed: %v", err)
		}
	}()

	// Port 443 handles SSL-terminated proxying
	go func() {
		server := &http.Server{
			Addr:      "0.0.0.0:443",
			Handler:   proxyHandler,
			TLSConfig: certManager.TLSConfig(),
		}
		log.Println("[proxy] HTTPS listener on port 443")
		if err := server.ListenAndServeTLS("", ""); err != nil {
			log.Fatalf("HTTPS proxy server failed: %v", err)
		}
	}()

	// Port LISTEN_PORT handles plain HTTP proxying
	server := &http.Server{
		Addr:    "0.0.0.0:" + listenPort,
		Handler: proxyHandler,
	}
	log.Printf("[proxy] HTTP listener on port %s", listenPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("HTTP proxy server failed: %v", err)
	}
}
