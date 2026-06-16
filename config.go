package main

import (
	"log"
	"sort"
)

func reloadConfigCache() {
	routeMu.Lock()
	defer routeMu.Unlock()

	var dbRoutes []ProxyRoute
	if err := db.Find(&dbRoutes).Error; err != nil {
		log.Printf("[config] failed to load routes from DB: %v", err)
		return
	}

	routeConfigDB = make(map[string]ProxyRoute)
	routeKeys = make([]string, 0, len(dbRoutes))

	for _, r := range dbRoutes {
		routeConfigDB[r.Domain] = r
		routeKeys = append(routeKeys, r.Domain)
	}

	// Pre-sort keys by descending length so the most-specific prefix wins
	sort.Slice(routeKeys, func(i, j int) bool {
		return len(routeKeys[i]) > len(routeKeys[j])
	})

	log.Printf("[config] cache reloaded: %d routes active", len(routeConfigDB))
}
