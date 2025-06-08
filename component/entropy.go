package component

import (
	"fmt"
	"os"

	"github.com/friedelschoen/st8/notify"
	"github.com/friedelschoen/st8/proto"
)

func entropyAvailable(args map[string]string, events *proto.EventHandlers) (Component, error) {
	return func(block *proto.Block, not *notify.Notification) error {
		data, err := os.ReadFile("/proc/sys/kernel/random/entropy_avail")
		if err != nil {
			return fmt.Errorf("unable to get entropy: %w", err)
		}
		block.Text = string(data)
		return nil
	}, nil
}

func init() {
	Install("entropy", entropyAvailable)
}
