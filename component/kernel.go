package component

import (
	"fmt"
	"syscall"

	"github.com/friedelschoen/st8/notify"
	"github.com/friedelschoen/st8/proto"
)

func kernelRelease(args map[string]string, events *proto.EventHandlers) (Component, error) {
	var res syscall.Utsname
	if err := syscall.Uname(&res); err != nil {
		return nil, fmt.Errorf("unable to get uname: %w", err)
	}
	bytes := make([]byte, len(res.Release))
	for i, chr := range res.Release {
		bytes[i] = byte(chr)
	}
	return func(block *proto.Block, not *notify.Notification) error {
		block.Text = string(bytes)
		return nil
	}, nil
}

func init() {
	Install("kernel_release", kernelRelease)
}
