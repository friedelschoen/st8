package component

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/friedelschoen/st8/notify"
)

var units = []string{"B", "kB", "MB", "GB"}

func getMem(key string) (uint64, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		name := strings.TrimSuffix(fields[0], ":")
		if key != name {
			continue
		}
		value, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return 0, err
		}
		if len(fields) > 2 {
			for _, unit := range units {
				if fields[2] == unit {
					break
				}
				value *= 1024
			}
		}
		return value, nil
	}
	return 0, scanner.Err()
}

func RamFree(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		avail, err := getMem("MemAvailable")
		if err != nil {
			return err
		}
		block.Text = fmtHuman(avail)
		return nil
	}, nil
}

func RamUsed(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		avail, err := getMem("MemAvailable")
		if err != nil {
			return err
		}
		total, err := getMem("MemTotal")
		if err != nil {
			return err
		}
		block.Text = fmtHuman(total - avail)
		return nil
	}, nil
}

func RamTotal(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		total, err := getMem("MemTotal")
		if err != nil {
			return err
		}
		block.Text = fmtHuman(total)
		return nil
	}, nil
}

func RamPercentage(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		avail, err := getMem("MemAvailable")
		if err != nil {
			return err
		}
		total, err := getMem("MemTotal")
		if err != nil {
			return err
		}
		block.Text = fmt.Sprintf("%.0f", 100-(float64(avail)/float64(total))*100)
		return nil
	}, nil
}
