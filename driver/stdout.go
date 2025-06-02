package driver

import (
	"fmt"

	"github.com/friedelschoen/st8/component"
)

type stdoutDriver struct{}

func init() {
	Drivers["stdout"] = stdoutDriver{}
}

func (stdoutDriver) SetText(line []component.Block) error {
	for i, block := range line {
		if i > 0 {
			fmt.Print(" | ")
		}
		fmt.Print(block.Text)
	}
	fmt.Println()
	return nil
}

func (stdoutDriver) Init() error {
	return nil
}

func (stdoutDriver) Close() error {
	return nil
}
