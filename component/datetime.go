package component

import (
	"fmt"
	"time"

	"github.com/friedelschoen/st8/notify"
	"github.com/ncruces/go-strftime"
)

func datetime(args map[string]string, events *EventHandlers) (Component, error) {
	datefmt, ok := args["datefmt"]
	if !ok {
		return nil, fmt.Errorf("missing argument: datefmt")
	}
	return func(block *Block, not *notify.Notification) error {
		block.Text = strftime.Format(datefmt, time.Now())
		return nil
	}, nil
}

func init() {
	Install("datetime", datetime)
}
