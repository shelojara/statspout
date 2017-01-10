package statspout

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
	Timestamp time.Time	`json:"@timestamp"`

	// associated container of this stats.
	ID string `json:"id"`
	Name string `json:"name"`

	// CPU usage percent.
	CpuPercent float64 `json:"cpu_percent"`

	// Memory usage in bytes.
	MemoryUsage uint64 `json:"mem_usage"`

	// Memory usage percent.
	MemoryPercent float64 `json:"mem_percent"`
}

func (stats *Stats) String() string {
	return fmt.Sprintf("[%s] {%s} CPU: %f%%, MEM: %f%% [%d Bytes]",
		stats.ID[:12], stats.Timestamp,
		stats.CpuPercent, stats.MemoryPercent, stats.MemoryUsage)
}
