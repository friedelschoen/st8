package component

import (
	"fmt"
	"os"

	"github.com/friedelschoen/st8/notify"
)

func EntropyAvailable(_ string, _ *notify.Notification) (string, error) {
	data, err := os.ReadFile("/proc/sys/kernel/random/entropy_avail")
	if err != nil {
		return "", fmt.Errorf("unable to get entropy: %w", err)
	}
	return string(data), nil
}
