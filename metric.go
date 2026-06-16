package main

import (
	"log"
	"time"
)

func metricWorker() {
	for {
		now := time.Now()
		showAt := now.Truncate(time.Minute)

		var existing ProxyMetric
		if err := db.Where("show_at = ?", showAt).First(&existing).Error; err != nil {
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

	var m ProxyMetric
	if err := db.Where("show_at = ?", showAt).First(&m).Error; err != nil {
		db.Create(&ProxyMetric{ShowAt: showAt, RequestVolume: 1, RequestLatency: float64(latencyMs)})
		return
	}

	newVolume := m.RequestVolume + 1
	newLatency := (m.RequestLatency*float64(m.RequestVolume) + float64(latencyMs)) / float64(newVolume)

	db.Model(&m).Updates(map[string]interface{}{
		"request_volume":  newVolume,
		"request_latency": newLatency,
	})
}

func startMetricWorker() {
	go metricWorker()
	log.Println("[metrics] worker started")
}
