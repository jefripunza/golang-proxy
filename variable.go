package main

import "sync"

var (
	listenPort   string
	sessionToken string

	// Mutex to protect routing cache
	routeMu sync.RWMutex

	// routeConfigDB maps "host/prefix" → ProxyRoute
	routeConfigDB map[string]ProxyRoute

	// sorted keys (longest first) for prefix-match priority
	routeKeys []string
)
