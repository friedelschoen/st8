package component

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/mem"
)

func SwapFree(_ string) (string, error) {
	v, err := mem.SwapMemory()
	if err != nil {
		return "", err
	}
	return fmtHuman(v.Free), nil
}

func SwapUsed(_ string) (string, error) {
	v, err := mem.SwapMemory()
	if err != nil {
		return "", err
	}
	used := v.Total - v.Free
	return fmtHuman(used), nil
}

func SwapTotal(_ string) (string, error) {
	v, err := mem.SwapMemory()
	if err != nil {
		return "", err
	}
	return fmtHuman(v.Total), nil
}

func SwapPercentage(_ string) (string, error) {
	v, err := mem.SwapMemory()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", int(v.UsedPercent)), nil
}
