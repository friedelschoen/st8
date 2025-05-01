package component

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/disk"
)

func DiskFree(path string) (string, error) {
	usage, err := disk.Usage(path)
	if err != nil {
		return "", err
	}
	return fmtHuman(usage.Free), nil
}

func DiskUsed(path string) (string, error) {
	usage, err := disk.Usage(path)
	if err != nil {
		return "", err
	}
	return fmtHuman(usage.Used), nil
}

func DiskTotal(path string) (string, error) {
	usage, err := disk.Usage(path)
	if err != nil {
		return "", err
	}
	return fmtHuman(usage.Total), nil
}

func DiskPercentage(path string) (string, error) {
	usage, err := disk.Usage(path)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", int(usage.UsedPercent)), nil
}
