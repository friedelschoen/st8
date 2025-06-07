package component

import (
	"time"

	"github.com/friedelschoen/st8/notify"
	"github.com/ncruces/go-strftime"
)

func Datetime(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		block.Text = strftime.Format(args["datefmt"], time.Now())
		return nil
	}, nil
}
