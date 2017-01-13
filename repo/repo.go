package repo

import (
	"github.com/mijara/statspout/stats"
)

// Repository Interface, used to define any repository.
// A repository is any service that can have data pushed to it.
// The close method is provided in case the service needs it.
type Interface interface {
	// Creates a new instance of the repository.
	// v holds the set of options for this new instance.
	Create(v interface{}) (Interface, error)

	// Push container stats to this service.
	// The repository should return an error if it's not capable of pushing the stats.
	Push(stats *stats.Stats) error

	// Close the service.
	Close()

	// Canonical name of this repository, used to identify it in the command line flags.
	Name() string
}
