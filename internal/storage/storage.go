package storage

import (
    "time"
    "message_queue_metrics/internal/monitor"
)

type Storage interface {
    SaveMetrics(metrics *monitor.Metrics) error
    GetMetrics(startTime, endTime time.Time) ([]*monitor.Metrics, error)
    Close() error
}