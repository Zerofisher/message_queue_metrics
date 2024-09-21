package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"message_queue_metrics/internal/storage"
)

type Handler struct {
	storage storage.Storage
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) GetMetrics(c *gin.Context) {
	// 获取查询参数
	startTimeStr := c.DefaultQuery("start_time", "")
	endTimeStr := c.DefaultQuery("end_time", "")

	// 解析时间
	var startTime, endTime time.Time
	var err error

	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_time format"})
			return
		}
	} else {
		startTime = time.Now().Add(-1 * time.Hour) // 默认查询过去一小时的数据
	}

	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_time format"})
			return
		}
	} else {
		endTime = time.Now()
	}

	// 查询数据
	metrics, err := h.storage.GetMetrics(startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve metrics"})
		return
	}

	c.JSON(http.StatusOK, metrics)
}
