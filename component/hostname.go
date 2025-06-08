package component

import (
	"os"

	"github.com/friedelschoen/st8/notify"
	"github.com/friedelschoen/st8/proto"
)

func hostname(args map[string]string, events *proto.EventHandlers) (Component, error) {
	host, err := os.Hostname()
	return func(block *proto.Block, not *notify.Notification) error {
		block.Text = host
		return err
	}, nil
}

func init() {
	Install("hostname", hostname)
}
