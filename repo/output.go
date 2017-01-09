package repo

import (
	"github.com/mijara/statspout/data"
)

type Interface interface {
	Push(stats *statspout.Stats)
}
