package driver

import (
	"io"

	"github.com/friedelschoen/st8/component"
)

type Driver interface {
	io.Closer

	Init() error
	SetText([]component.Block) error
}

var Drivers = map[string]Driver{}
