package component

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/mem"
)

func RamFree(_ string) (string, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}
	return fmtHuman(v.Available), nil
}

func RamUsed(_ string) (string, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}
	used := v.Total - v.Available
	return fmtHuman(used), nil
}

func RamTotal(_ string) (string, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}
	return fmtHuman(v.Total), nil
}

func RamPercentage(_ string) (string, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", int(v.UsedPercent)), nil
}
