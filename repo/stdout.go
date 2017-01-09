package repo

import (
	"fmt"

	"github.com/mijara/statspout/data"
)

type Stdout struct {

}

func (out Stdout) Push(stats *statspout.Stats) {
	fmt.Println(stats)
}
