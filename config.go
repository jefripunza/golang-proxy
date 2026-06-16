package main

import (
	"encoding/json"
	"log"
	"os"
	"sort"
)

func setConfig(configFile string) {
	// 1. Check if DB has any routes. If not, seed from config.json
	var count int64
	db.Model(&ProxyRoute{}).Count(&count)
	if count == 0 {
		data, err := os.ReadFile(configFile)
		if err == nil {
			var legacyConfig map[string]RouteConfig
			if err := json.Unmarshal(data, &legacyConfig); err == nil {
				for domain, cfg := range legacyConfig {
					route := ProxyRoute{
						Domain:     domain,
						SchemaType: "static",
						TargetURL:  cfg.To,
					}
					if cfg.BasicAuth != nil {
						route.UseBasicAuth = true
						route.BasicAuthUsername = cfg.BasicAuth.Username
						route.BasicAuthPassword = cfg.BasicAuth.Password
					}
					if err := db.Create(&route).Error; err != nil {
						log.Printf("[config] failed to seed legacy route %s: %v", domain, err)
					}
				}
				log.Printf("[config] seeded %d legacy routes from %s to SQLite", len(legacyConfig), configFile)
			}
		}
	}

	reloadConfigCache()
}

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
