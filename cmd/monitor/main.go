package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"message_queue_metrics/internal/api"
	"message_queue_metrics/internal/monitor"
	"message_queue_metrics/internal/storage"
)

func main() {
	brokers := []string{"127.0.0.1:9092"}
	topic := "test-topic"

	kafkaMonitor := monitor.NewKafkaMonitor(brokers, topic)

	// 初始化 BuntDB 存储
	buntStorage, err := storage.NewBuntDBStorage("metrics.db")
	if err != nil {
		log.Fatalf("Failed to initialize BuntDB storage: %v", err)
	}
	defer buntStorage.Close()

	// 初始化 Gin 路由
	router := gin.Default()
	handler := api.NewHandler(buntStorage)
	router.GET("/metrics", handler.GetMetrics)

	// 启动 API 服务器
	go func() {
		if err := router.Run(":8080"); err != nil {
			log.Fatalf("Failed to start Gin server: %v", err)
		}
	}()

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

			// 将指标保存到 BuntDB
			err = buntStorage.SaveMetrics(metrics)
			if err != nil {
				log.Printf("Error saving metrics to BuntDB: %v", err)
			} else {
				log.Printf("Metrics saved to BuntDB: %+v", metrics)
			}
		}
	}
}
