package monitor

type Metrics struct {
	MessageCount int
	PartitionLag map[int]int64
}
