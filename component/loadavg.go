package component

import (
	"fmt"

	"github.com/friedelschoen/st8/notify"
	"github.com/shirou/gopsutil/v3/load"
)

func LoadAverage(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	stat, err := load.Avg()
	if err != nil {
		return fmt.Errorf("unable to get average load: %w", err)
	}
	switch args["minutes"] {
	case "":
		block.Text = fmt.Sprintf("%.2f %.2f %.2f", stat.Load1, stat.Load5, stat.Load15)
		return nil
	case "1":
		block.Text = fmt.Sprintf("%.2f", stat.Load1)
		return nil
	case "5":
		block.Text = fmt.Sprintf("%.2f", stat.Load5)
		return nil
	case "15":
		block.Text = fmt.Sprintf("%.2f", stat.Load15)
		return nil
	default:
		return fmt.Errorf("unable to get average load: period must be either 1, 5 or 15 (minutes)")
	}
}
