package repo

import (
	"github.com/mijara/statspout/data"
)

// Repository Interface, used to define any repository.
// A repository is any service that can have data pushed to it.
// The close method is provided in case the service needs it.
type Interface interface {
	// Push container stats to this service.
	Push(stats *statspout.Stats) error

	// Close the service.
	Close()
}
