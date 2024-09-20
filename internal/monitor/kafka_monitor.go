package monitor

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaMonitor struct {
	brokers []string
	topic   string
}

func NewKafkaMonitor(brokers []string, topic string) *KafkaMonitor {
	return &KafkaMonitor{
		brokers: brokers,
		topic:   topic,
	}
}

func (km *KafkaMonitor) CollectMetrics() (*Metrics, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", km.brokers[0], km.topic, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to dial leader: %v", err)
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return nil, fmt.Errorf("failed to read partitions: %v", err)
	}

	metrics := &Metrics{
		PartitionLag: make(map[int]int64),
	}

	for _, p := range partitions {
		high, err := conn.ReadLastOffset()
		if err != nil {
			return nil, fmt.Errorf("failed to read last offset: %v", err)
		}

		low, err := conn.ReadFirstOffset()
		if err != nil {
			return nil, fmt.Errorf("failed to read first offset: %v", err)
		}

		metrics.MessageCount += int(high - low)
		metrics.PartitionLag[p.ID] = high - low
	}

	return metrics, nil
}
