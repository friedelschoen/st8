package component

import (
	"fmt"
	"syscall"

	"github.com/friedelschoen/st8/notify"
)

func kernelRelease(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		var res syscall.Utsname
		if err := syscall.Uname(&res); err != nil {
			return fmt.Errorf("unable to get uname: %w", err)
		}
		bytes := make([]byte, len(res.Release))
		for i, chr := range res.Release {
			bytes[i] = byte(chr)
		}
		block.Text = string(bytes)
		return nil
	}, nil
}

func init() {
	Install("kernel_release", kernelRelease)
}
