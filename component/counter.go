package component

import (
	"strconv"

	"github.com/friedelschoen/st8/notify"
)

func counter(args map[string]string, events *EventHandlers) (Component, error) {
	count := 0

	events.OnClick = func(ClickEvent) {
		count++
	}

	return func(block *Block, not *notify.Notification) error {
		block.Text = strconv.Itoa(count)
		return nil
	}, nil
}

func init() {
	Install("counter", counter)
}
