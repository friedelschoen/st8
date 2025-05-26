package driver

import "io"

type Driver interface {
	io.Closer

	Init() error
	SetText(string) error
}

var Drivers = map[string]Driver{}
