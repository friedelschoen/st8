package component

import (
	"fmt"

	"github.com/friedelschoen/st8/notify"
	"github.com/shirou/gopsutil/v3/mem"
)

func RamFree(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	v, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	block.Text = fmtHuman(v.Available)
	return nil
}

func RamUsed(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	v, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	used := v.Total - v.Available
	block.Text = fmtHuman(used)
	return nil
}

func RamTotal(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	v, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	block.Text = fmtHuman(v.Total)
	return nil
}

func RamPercentage(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	v, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	block.Text = fmt.Sprintf("%d", int(v.UsedPercent))
	return nil
}
