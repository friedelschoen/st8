package driver

import (
	"io"

	"github.com/friedelschoen/st8/proto"
)

type Driver interface {
	io.Closer

	Init(update chan<- struct{}) error
	SetText([]proto.Block) error
}

var Drivers = map[string]Driver{}
