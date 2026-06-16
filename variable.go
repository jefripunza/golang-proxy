package main

import (
	"encoding/json"
	"sync"
)

var (
	listenPort   string
	sessionToken string

	// Mutex to protect routing cache
	routeMu sync.RWMutex

	// routeConfigDB maps "host/prefix" → ProxyRoute
	routeConfigDB map[string]ProxyRoute

	// sorted keys (longest first) for prefix-match priority
	routeKeys []string

	// SSE broker for real-time log streaming
	sseBroker = &logBroker{
		clients:    make(map[chan []byte]string),
		register:   make(chan subscription),
		unregister: make(chan chan []byte),
		broadcast:  make(chan ProxyLog, 256),
	}
)

type subscription struct {
	ch     chan []byte
	domain string
}

type logBroker struct {
	mu         sync.Mutex
	clients    map[chan []byte]string
	register   chan subscription
	unregister chan chan []byte
	broadcast  chan ProxyLog
}

func (b *logBroker) run() {
	for {
		select {
		case sub := <-b.register:
			b.mu.Lock()
			b.clients[sub.ch] = sub.domain
			b.mu.Unlock()
		case ch := <-b.unregister:
			b.mu.Lock()
			delete(b.clients, ch)
			close(ch)
			b.mu.Unlock()
		case log := <-b.broadcast:
			data, err := json.Marshal(log)
			if err != nil {
				continue
			}
			b.mu.Lock()
			for ch, domain := range b.clients {
				if domain == "" || domain == log.Domain {
					select {
					case ch <- data:
					default:
					}
				}
			}
			b.mu.Unlock()
		}
	}
}
