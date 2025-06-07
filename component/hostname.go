package component

import (
	"os"

	"github.com/friedelschoen/st8/notify"
)

func Hostname(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		var err error
		block.Text, err = os.Hostname()
		return err
	}, nil
}
