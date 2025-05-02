package component

import (
	"fmt"

	"github.com/friedelschoen/st8/notify"
	"github.com/shirou/gopsutil/v3/load"
)

func LoadAverage(min string, _ *notify.Notification, _ *any) (string, error) {
	stat, err := load.Avg()
	if err != nil {
		return "", fmt.Errorf("unable to get average load: %w", err)
	}
	switch min {
	case "":
		return fmt.Sprintf("%.2f %.2f %.2f", stat.Load1, stat.Load5, stat.Load15), nil
	case "1":
		return fmt.Sprintf("%.2f", stat.Load1), nil
	case "5":
		return fmt.Sprintf("%.2f", stat.Load5), nil
	case "15":
		return fmt.Sprintf("%.2f", stat.Load15), nil
	default:
		return "", fmt.Errorf("unable to get average load: period must be either 1, 5 or 15 (minutes)")
	}
}
