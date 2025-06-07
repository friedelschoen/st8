package component

import (
	"fmt"
	"os"

	"github.com/friedelschoen/st8/notify"
)

func EntropyAvailable(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		data, err := os.ReadFile("/proc/sys/kernel/random/entropy_avail")
		if err != nil {
			return fmt.Errorf("unable to get entropy: %w", err)
		}
		block.Text = string(data)
		return nil
	}, nil
}
