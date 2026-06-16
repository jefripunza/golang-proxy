package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	// Load .env (silent if file absent – env vars already set are respected)
	_ = godotenv.Load()

	listenPort = os.Getenv("LISTEN_PORT")
	if listenPort == "" {
		listenPort = "8080"
	}

	// Generate a random 32-byte session token
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		log.Fatalf("cannot generate session token: %v", err)
	}
	sessionToken = hex.EncodeToString(raw)
}
