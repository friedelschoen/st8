package component

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

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

type cpucache struct {
	lastCPUTimes *CPUStat
	lastTime     time.Time
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

func CPUPercentage(block *Block, args map[string]string, not *notify.Notification, cacheptr *any) error {
	var cache cpucache
	if *cacheptr != nil {
		cache = (*cacheptr).(cpucache)
	}

	curTimes, err := times()
	if err != nil || curTimes == nil {
		return fmt.Errorf("unable to get CPU times: %w", err)
	}

	if cache.lastCPUTimes == nil {
		cache.lastCPUTimes = curTimes
		cache.lastTime = time.Now()
		*cacheptr = cache
		block.Text = "0"
		return nil // first call
	}

	last := cache.lastCPUTimes
	curr := curTimes

	totalDelta := totalCPU(curr) - totalCPU(last)
	idleDelta := curr.Idle - last.Idle

	if totalDelta <= 0 {
		block.Text = "0"
		return nil
	}

	usage := 100.0 * (1.0 - idleDelta/totalDelta)

	cache.lastCPUTimes = curTimes
	cache.lastTime = time.Now()
	*cacheptr = cache

	block.Text = fmt.Sprintf("%.0f", usage)
	return nil
}

func totalCPU(t *CPUStat) float64 {
	return t.User + t.System + t.Idle + t.Nice + t.Iowait + t.Irq + t.Softirq + t.Steal + t.Guest + t.GuestNice
}
