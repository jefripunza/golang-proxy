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

	autoSSL := os.Getenv("AUTO_SSL") == "true"

	if autoSSL {
		// Auto SSL Let's Encrypt setup
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: hostPolicy,
			Cache:      autocert.DirCache("certs"),
		}

		// Port 80 handles HTTP-01 challenges and redirects others to HTTPS
		go func() {
			log.Println("[proxy] HTTP server redirecting to HTTPS on port 80...")
			err := http.ListenAndServe("0.0.0.0:80", certManager.HTTPHandler(nil))
			if err != nil {
				log.Fatalf("HTTP challenge listener failed: %v", err)
			}
		}()

		// Port 443 handles the actual proxying securely
		server := &http.Server{
			Addr:      "0.0.0.0:443",
			Handler:   proxyHandler,
			TLSConfig: certManager.TLSConfig(),
		}

		log.Println("[proxy] listening securely on port 443 (HTTPS)")
		if err := server.ListenAndServeTLS("", ""); err != nil {
			log.Fatalf("HTTPS proxy server failed: %v", err)
		}
	} else {
		// Normal HTTP listener
		server := &http.Server{
			Addr:    "0.0.0.0:" + listenPort,
			Handler: proxyHandler,
		}

		log.Printf("[proxy] listening on port %s (HTTP)", listenPort)
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("HTTP proxy server failed: %v", err)
		}
	}
}
