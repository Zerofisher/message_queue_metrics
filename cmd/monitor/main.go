package main

import (
	"log"
	"time"

	"message_queue_metrics/internal/monitor"
)

func main() {
	brokers := []string{"127.0.0.1:9092"}
	topic := "test-topic"

	kafkaMonitor := monitor.NewKafkaMonitor(brokers, topic)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			metrics, err := kafkaMonitor.CollectMetrics()
			if err != nil {
				log.Printf("Error collecting metrics: %v", err)
				continue
			}
			log.Printf("Metrics: %+v", metrics)
		}
	}
}
