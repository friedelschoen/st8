package component

import (
	"fmt"

	"github.com/friedelschoen/st8/notify"
)

func swapFree(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		avail, err := getMem("SwapFree")
		if err != nil {
			return err
		}
		block.Text = fmtHuman(avail)
		return nil
	}, nil
}

func swapUsed(args map[string]string, events *EventHandlers) (Component, error) {
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

func swapTotal(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		total, err := getMem("SwapTotal")
		if err != nil {
			return err
		}
		block.Text = fmtHuman(total)
		return nil
	}, nil
}

func swapPercentage(args map[string]string, events *EventHandlers) (Component, error) {
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

func init() {
	Install("swap_free", swapFree)
	Install("swap_perc", swapPercentage)
	Install("swap_total", swapTotal)
	Install("swap_used", swapUsed)
}
