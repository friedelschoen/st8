package component

import (
	"fmt"
	"time"

	"github.com/friedelschoen/st8/notify"
	"github.com/shirou/gopsutil/v3/cpu"
)

type cpucache struct {
	lastCPUTimes []cpu.TimesStat
	lastTime     time.Time
}

// cpuFreq returns the average current frequency of all CPUs in Hz as a formatted string.
func CPUFrequency(_ string, _ *notify.Notification, _ *any) (string, error) {
	freqs, err := cpu.Info()
	if err != nil || len(freqs) == 0 {
		return "", fmt.Errorf("unable to get CPU frequency: %w", err)
	}

	var sum float64
	for _, f := range freqs {
		sum += f.Mhz
	}
	avgFreqMHz := sum / float64(len(freqs))
	return fmt.Sprintf("%.0f MHz", avgFreqMHz), nil
}

func CPUPercentage(_ string, _ *notify.Notification, cacheptr *any) (string, error) {
	var cache cpucache
	if *cacheptr != nil {
		cache = (*cacheptr).(cpucache)
	}

	curTimes, err := cpu.Times(false)
	if err != nil || len(curTimes) == 0 {
		return "", fmt.Errorf("unable to get CPU times: %w", err)
	}

	if len(cache.lastCPUTimes) == 0 {
		cache.lastCPUTimes = curTimes
		cache.lastTime = time.Now()
		*cacheptr = cache
		return "0", nil // first call
	}

	last := cache.lastCPUTimes[0]
	curr := curTimes[0]

	totalDelta := totalCPU(curr) - totalCPU(last)
	idleDelta := curr.Idle - last.Idle

	if totalDelta <= 0 {
		return "0", nil
	}

	usage := 100.0 * (1.0 - idleDelta/totalDelta)

	cache.lastCPUTimes = curTimes
	cache.lastTime = time.Now()
	*cacheptr = cache

	return fmt.Sprintf("%.0f", usage), nil
}

func totalCPU(t cpu.TimesStat) float64 {
	return t.User + t.System + t.Idle + t.Nice + t.Iowait + t.Irq + t.Softirq + t.Steal + t.Guest + t.GuestNice
}
