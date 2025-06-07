package component

import (
	"fmt"

	"github.com/friedelschoen/st8/notify"
)

func SwapFree(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		avail, err := getMem("SwapFree")
		if err != nil {
			return err
		}
		block.Text = fmtHuman(avail)
		return nil
	}, nil
}

func SwapUsed(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		avail, err := getMem("SwapFree")
		if err != nil {
			return err
		}
		total, err := getMem("SwapTotal")
		if err != nil {
			return err
		}
		block.Text = fmtHuman(total - avail)
		return nil
	}, nil
}

func SwapTotal(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		total, err := getMem("SwapTotal")
		if err != nil {
			return err
		}
		block.Text = fmtHuman(total)
		return nil
	}, nil
}

func SwapPercentage(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		avail, err := getMem("SwapFree")
		if err != nil {
			return err
		}
		total, err := getMem("SwapTotal")
		if err != nil {
			return err
		}
		block.Text = fmt.Sprintf("%.0f", 100-(float64(avail)/float64(total))*100)
		return nil
	}, nil
}
