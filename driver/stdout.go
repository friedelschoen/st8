package driver

import "fmt"

type stdoutDriver struct{}

func init() {
	Drivers["stdout"] = stdoutDriver{}
}

func (stdoutDriver) SetText(line string) error {
	_, err := fmt.Println(line)
	return err
}

func (stdoutDriver) Init() error {
	return nil
}

func (stdoutDriver) Close() error {
	return nil
}
