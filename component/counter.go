package component

import (
	"strconv"

	"github.com/friedelschoen/st8/notify"
	"github.com/friedelschoen/st8/proto"
)

func counter(args map[string]string, events *proto.EventHandlers) (Component, error) {
	count := 0

	events.OnClick = func(proto.ClickEvent) {
		count++
	}

	return func(block *proto.Block, not *notify.Notification) error {
		block.Text = strconv.Itoa(count)
		return nil
	}, nil
}

func init() {
	Install("counter", counter)
}
