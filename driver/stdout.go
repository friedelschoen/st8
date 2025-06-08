package driver

import (
	"fmt"

	"github.com/friedelschoen/st8/proto"
)

type stdoutDriver struct{}

func init() {
	Drivers["stdout"] = stdoutDriver{}
}

func (stdoutDriver) SetText(line []proto.Block) error {
	for i, block := range line {
		if i > 0 {
			fmt.Print(" | ")
		}
		fmt.Print(block.Text)
	}
	fmt.Println()
	return nil
}

func (stdoutDriver) Init(chan<- struct{}) error {
	return nil
}

func (stdoutDriver) Close() error {
	return nil
}
