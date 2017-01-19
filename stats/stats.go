package stats

import (
	"fmt"
	"time"
)

// Standard project container to not relay on a specific docker client implementation.
type Container struct {
	ID    string
	Names []string
	Image string
}

// Standard project stats to not relay on a specific docker client implementation.
type Stats struct {
	// Timestamp of this stats.
	Timestamp time.Time    `json:"@timestamp"`

	// associated container of this stats.
	Name string `json:"name"`

	// CPU usage percent.
	CpuPercent float64 `json:"cpu_percent"`

	// Memory usage in bytes.
	MemoryUsage uint64 `json:"mem_usage"`

	// Memory usage percent.
	MemoryPercent float64 `json:"mem_percent"`

	// Transmit and Receive network stats, in bytes.
	TxBytesTotal uint32 `json:"tx_bytes"`
	RxBytesTotal uint32 `json:"rx_bytes"`
}

// Prints stats in a nice format.
func (stats *Stats) String() string {
	return fmt.Sprintf("[%s] {%s} CPU: %f%%, MEM: %f%% [%d B] Tx/Rx: %d/%d",
		stats.Name, stats.Timestamp.Format("02 Jan 06 15:04:05 MST"),
		stats.CpuPercent, stats.MemoryPercent, stats.MemoryUsage,
		stats.TxBytesTotal, stats.RxBytesTotal)
}
