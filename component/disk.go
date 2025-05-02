package component

import (
	"fmt"

	"github.com/friedelschoen/st8/notify"
	"github.com/shirou/gopsutil/v3/disk"
)

func DiskFree(path string, _ *notify.Notification) (string, error) {
	usage, err := disk.Usage(path)
	if err != nil {
		return "", err
	}
	return fmtHuman(usage.Free), nil
}

func DiskUsed(path string, _ *notify.Notification) (string, error) {
	usage, err := disk.Usage(path)
	if err != nil {
		return "", err
	}
	return fmtHuman(usage.Used), nil
}

func DiskTotal(path string, _ *notify.Notification) (string, error) {
	usage, err := disk.Usage(path)
	if err != nil {
		return "", err
	}
	return fmtHuman(usage.Total), nil
}

func DiskPercentage(path string, _ *notify.Notification) (string, error) {
	usage, err := disk.Usage(path)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", int(usage.UsedPercent)), nil
}
