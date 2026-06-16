package main

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	os.MkdirAll("database", 0755)

	var err error
	db, err = gorm.Open(sqlite.Open("database/proxy.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&ProxyRoute{}, &ProxyLog{}, &BlockedIP{}, &ProxyMetric{}, &RateLimitRecord{}, &AutoClearSetting{})
	// Ensure default auto-clear setting exists
	var count int64
	db.Model(&AutoClearSetting{}).Count(&count)
	if count == 0 {
		db.Create(&AutoClearSetting{Interval: 0}) // 0 = never
	}
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	log.Println("[db] successfully connected and migrated SQLite database")
}
