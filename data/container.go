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
	Timestamp time.Time

	// associated container of this stats.
	Container *Container

	// CPU usage percent.
	CpuPercent float64

	// Memory usage in bytes.
	MemoryUsage uint64

	// Memory usage percent.
	MemoryPercent float64
}

func (stats *Stats) String() string {
	return fmt.Sprintf("[%s] {%s} CPU: %f%%, MEM: %f%% [%d Bytes]",
		stats.Container.ID[:12], stats.Timestamp,
		stats.CpuPercent, stats.MemoryPercent, stats.MemoryUsage)
}
