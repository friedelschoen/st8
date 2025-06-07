package component

import (
	"os"

	"github.com/friedelschoen/st8/notify"
)

func hostname(args map[string]string, events *EventHandlers) (Component, error) {
	host, err := os.Hostname()
	return func(block *Block, not *notify.Notification) error {
		block.Text = host
		return err
	}, nil
}

func init() {
	Install("hostname", hostname)
}
