package component

import (
	"fmt"
	"time"

	"github.com/friedelschoen/st8/notify"
	"github.com/shirou/gopsutil/v3/host"
)

func Uptime(_ string, _ *notify.Notification) (string, error) {
	seconds, err := host.Uptime()
	if err != nil {
		return "", fmt.Errorf("unable to get uptime: %w", err)
	}
	return (time.Duration(seconds) * time.Second).String(), nil
}
