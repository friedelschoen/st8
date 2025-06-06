package component

import (
	"fmt"

	"github.com/friedelschoen/st8/notify"
)

func SwapFree(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	avail, err := getMem("SwapFree")
	if err != nil {
		return err
	}
	block.Text = fmtHuman(avail)
	return nil
}

func SwapUsed(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
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
}

func SwapTotal(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	total, err := getMem("SwapTotal")
	if err != nil {
		return err
	}
	block.Text = fmtHuman(total)
	return nil
}

func SwapPercentage(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
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
}
