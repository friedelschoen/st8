package component

import (
	"fmt"

	"github.com/friedelschoen/st8/notify"
	"github.com/shirou/gopsutil/v3/mem"
)

func SwapFree(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	v, err := mem.SwapMemory()
	if err != nil {
		return err
	}
	block.Text = fmtHuman(v.Free)
	return nil
}

func SwapUsed(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	v, err := mem.SwapMemory()
	if err != nil {
		return err
	}
	used := v.Total - v.Free
	block.Text = fmtHuman(used)
	return nil
}

func SwapTotal(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	v, err := mem.SwapMemory()
	if err != nil {
		return err
	}
	block.Text = fmtHuman(v.Total)
	return nil
}

func SwapPercentage(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	v, err := mem.SwapMemory()
	if err != nil {
		return err
	}
	block.Text = fmt.Sprintf("%d", int(v.UsedPercent))
	return nil
}
