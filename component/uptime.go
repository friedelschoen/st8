package component

import (
	"fmt"
	"time"

	"github.com/friedelschoen/st8/notify"
	"github.com/shirou/gopsutil/v3/host"
)

func Uptime(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	seconds, err := host.Uptime()
	if err != nil {
		return fmt.Errorf("unable to get uptime: %w", err)
	}
	block.Text = (time.Duration(seconds) * time.Second).String()
	return nil
}
