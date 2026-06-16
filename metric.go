package main

import (
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func metricWorker() {
	for {
		now := time.Now()
		showAt := now.Truncate(time.Minute)

		var existing ProxyMetric
		if err := db.
		Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}). // silent mode to avoid noise
		Where("show_at = ?", showAt).First(&existing).Error; err != nil {
			db.Create(&ProxyMetric{ShowAt: showAt, RequestVolume: 0, RequestLatency: 0})
		}

		next := showAt.Add(time.Minute)
		sleepDuration := next.Sub(time.Now())
		if sleepDuration > 0 {
			time.Sleep(sleepDuration)
		}
	}
}

func recordRequestMetric(latencyMs int64) {
	showAt := time.Now().Truncate(time.Minute)
	lat := float64(latencyMs)

	var m ProxyMetric
	if err := db.Where("show_at = ?", showAt).First(&m).Error; err != nil {
		db.Create(&ProxyMetric{ShowAt: showAt, RequestVolume: 1, RequestLatency: lat, MaxLatency: lat, MinLatency: lat})
		return
	}

	newVolume := m.RequestVolume + 1
	newLatency := (m.RequestLatency*float64(m.RequestVolume) + lat) / float64(newVolume)
	maxLat := m.MaxLatency
	minLat := m.MinLatency
	if lat > maxLat {
		maxLat = lat
	}
	if m.RequestVolume == 0 || lat < minLat {
		minLat = lat
	}

	db.Model(&m).Updates(map[string]interface{}{
		"request_volume":  newVolume,
		"request_latency": newLatency,
		"max_latency":     maxLat,
		"min_latency":     minLat,
	})
}

func autoClearWorker() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		var setting AutoClearSetting
		if err := db.First(&setting).Error; err != nil || setting.Interval == 0 {
			continue
		}
		cutoff := time.Now().Add(-time.Duration(setting.Interval) * time.Hour)
		db.Where("timestamp < ?", cutoff).Delete(&ProxyLog{})
	}
}

func startAutoClearWorker() {
	go autoClearWorker()
	log.Println("[auto-clear] worker started")
}

func startMetricWorker() {
	go metricWorker()
	log.Println("[metrics] worker started")
}
