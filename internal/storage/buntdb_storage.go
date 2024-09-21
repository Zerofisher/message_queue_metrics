package storage

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tidwall/buntdb"
	"message_queue_metrics/internal/monitor"
)

type BuntDBStorage struct {
	db *buntdb.DB
}

var _ Storage = (*BuntDBStorage)(nil)

func NewBuntDBStorage(path string) (*BuntDBStorage, error) {
	db, err := buntdb.Open(path)
	if err != nil {
		return nil, err
	}
	return &BuntDBStorage{db: db}, nil
}

func (s *BuntDBStorage) SaveMetrics(metrics *monitor.Metrics) error {
	return s.db.Update(func(tx *buntdb.Tx) error {
		key := fmt.Sprintf("metrics:%d", time.Now().UnixNano())
		value, err := json.Marshal(metrics)
		if err != nil {
			return err
		}
		_, _, err = tx.Set(key, string(value), nil)
		return err
	})
}

func (s *BuntDBStorage) GetMetrics(startTime, endTime time.Time) ([]*monitor.Metrics, error) {
	var results []*monitor.Metrics
	err := s.db.View(func(tx *buntdb.Tx) error {
		return tx.AscendRange("",
			fmt.Sprintf("metrics:%d", startTime.UnixNano()),
			fmt.Sprintf("metrics:%d", endTime.UnixNano()),
			func(key, value string) bool {
				var m monitor.Metrics
				if err := json.Unmarshal([]byte(value), &m); err != nil {
					return false
				}
				results = append(results, &m)
				return true
			})
	})
	return results, err
}

func (s *BuntDBStorage) Close() error {
	return s.db.Close()
}
