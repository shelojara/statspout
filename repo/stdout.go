package repo

import (
	"fmt"

	"github.com/mijara/statspout/data"
)

type Stdout struct {

}

func NewStdout() *Stdout {
	return &Stdout{}
}

func (out *Stdout) Push(stats *statspout.Stats) error {
	fmt.Println(stats)
	return nil
}
