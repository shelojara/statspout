package common

import (
	"fmt"

	"github.com/mijara/statspout/data"
)

type Stdout struct {
}

func NewStdout() *Stdout {
	return &Stdout{}
}

func (stdout *Stdout) Push(stats *statspout.Stats) error {
	fmt.Println(stats)
	return nil
}

func (stdout *Stdout) Close() {

}
