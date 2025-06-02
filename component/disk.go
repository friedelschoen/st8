package component

import (
	"fmt"

	"github.com/friedelschoen/st8/notify"
	"github.com/shirou/gopsutil/v3/disk"
)

func DiskFree(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	usage, err := disk.Usage(args["path"])
	if err != nil {
		return err
	}
	block.Text = fmtHuman(usage.Free)
	return nil
}

func DiskUsed(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	usage, err := disk.Usage(args["path"])
	if err != nil {
		return err
	}
	block.Text = fmtHuman(usage.Used)
	return nil
}

func DiskTotal(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	usage, err := disk.Usage(args["path"])
	if err != nil {
		return err
	}
	block.Text = fmtHuman(usage.Total)
	return nil
}

func DiskPercentage(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	usage, err := disk.Usage(args["path"])
	if err != nil {
		return err
	}
	block.Text = fmt.Sprintf("%d", int(usage.UsedPercent))
	return nil
}
