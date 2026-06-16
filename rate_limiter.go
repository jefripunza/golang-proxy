package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

type rateBucket struct {
	count   int
	resetAt time.Time
}

type rateLimiter struct {
	mu      sync.Mutex
	buckets map[string]map[string]*rateBucket
}

var rl = &rateLimiter{
	buckets: make(map[string]map[string]*rateBucket),
}

func (r *rateLimiter) allow(routeKey, clientKey string, limit int, window time.Duration) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	if r.buckets[routeKey] == nil {
		r.buckets[routeKey] = make(map[string]*rateBucket)
	}

	b, ok := r.buckets[routeKey][clientKey]
	if !ok || now.After(b.resetAt) {
		r.buckets[routeKey][clientKey] = &rateBucket{count: 1, resetAt: now.Add(window)}
		return true
	}

	if b.count >= limit {
		return false
	}

	b.count++
	return true
}

func (r *rateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		r.mu.Lock()
		now := time.Now()
		for route, clients := range r.buckets {
			for key, b := range clients {
				if now.After(b.resetAt) {
					delete(clients, key)
				}
			}
			if len(clients) == 0 {
				delete(r.buckets, route)
			}
		}
		r.mu.Unlock()
	}
}

func unitToDuration(unit string) time.Duration {
	switch unit {
	case "second":
		return time.Second
	case "minute":
		return time.Minute
	case "hour":
		return time.Hour
	case "day":
		return 24 * time.Hour
	default:
		return time.Minute
	}
}

var rateLimitNamespace = uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

func getRateLimitKey(r *http.Request, route ProxyRoute, clientIP string) string {
	switch route.RateLimitMethod {
	case "ip":
		return clientIP
	case "header":
		val := r.Header.Get(route.RateLimitHeaderKey)
		if route.RateLimitHeaderValue != "" && val != route.RateLimitHeaderValue {
			return ""
		}
		if val == "" {
			return clientIP
		}
		return val
	default:
		// compact: serialize all headers to JSON, hash with UUIDv5 for a compact unique fingerprint
		headerMap := make(map[string]string)
		for k, vv := range r.Header {
			headerMap[k] = vv[0]
		}
		data, _ := json.Marshal(headerMap)
		return uuid.NewSHA1(rateLimitNamespace, data).String()
	}
}
