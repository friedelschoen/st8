package component

import (
	"fmt"
	"os"
)

func EntropyAvailable(_ string) (string, error) {
	data, err := os.ReadFile("/proc/sys/kernel/random/entropy_avail")
	if err != nil {
		return "", fmt.Errorf("unable to get entropy: %w", err)
	}
	return string(data), nil
}
