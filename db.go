package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("proxy.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&ProxyRoute{}, &ProxyLog{}, &BlockedIP{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	log.Println("[db] successfully connected and migrated SQLite database")
}
