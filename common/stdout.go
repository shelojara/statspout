package common

import (
	"fmt"

	"github.com/mijara/statspout/repo"
	"github.com/mijara/statspout/stats"
)

type Stdout struct {
}

func (*Stdout) Name() string {
	return "stdout"
}

func (*Stdout) Create(v interface{}) (repo.Interface, error) {
	return NewStdout(), nil
}

func (*Stdout) Clear(name string) {
}

func NewStdout() *Stdout {
	return &Stdout{}
}

func (stdout *Stdout) Push(s *stats.Stats) error {
	fmt.Println(s)
	return nil
}

func (stdout *Stdout) Close() {

}
