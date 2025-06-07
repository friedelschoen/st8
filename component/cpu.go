package component

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/friedelschoen/st8/notify"
)

type CPUStat struct {
	CPU       string
	User      float64
	System    float64
	Idle      float64
	Nice      float64
	Iowait    float64
	Irq       float64
	Softirq   float64
	Steal     float64
	Guest     float64
	GuestNice float64
}

func times() (*CPUStat, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if fields[0] != "cpu" {
			continue
		}
		var cpu CPUStat
		cpu.CPU = fields[0]
		dests := []*float64{
			&cpu.User,
			&cpu.Nice,
			&cpu.System,
			&cpu.Idle,
			&cpu.Iowait,
			&cpu.Irq,
			&cpu.Softirq,
			&cpu.Steal,
			&cpu.Guest,
			&cpu.GuestNice,
		}
		for i, dst := range dests {
			if *dst, err = strconv.ParseFloat(fields[i+1], 64); err != nil {
				return nil, err
			}
		}
		return &cpu, nil
	}
	return nil, scanner.Err()
}

func CPUPercentage(args map[string]string, events *EventHandlers) (Component, error) {
	var lastCPUTimes *CPUStat

	return func(block *Block, not *notify.Notification) error {
		curTimes, err := times()
		if err != nil || curTimes == nil {
			return fmt.Errorf("unable to get CPU times: %w", err)
		}

		if lastCPUTimes == nil {
			lastCPUTimes = curTimes
			block.Text = "0"
			return nil
		}

		curr := curTimes

		totalDelta := totalCPU(curr) - totalCPU(lastCPUTimes)
		idleDelta := curr.Idle - lastCPUTimes.Idle

		if totalDelta <= 0 {
			block.Text = "0"
			return nil
		}

		usage := 100.0 * (1.0 - idleDelta/totalDelta)

		lastCPUTimes = curTimes

		block.Text = fmt.Sprintf("%.0f", usage)
		return nil
	}, nil
}

func totalCPU(t *CPUStat) float64 {
	return t.User + t.System + t.Idle + t.Nice + t.Iowait + t.Irq + t.Softirq + t.Steal + t.Guest + t.GuestNice
}
