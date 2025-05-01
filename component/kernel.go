package component

import (
	"fmt"
	"syscall"
)

func KernelRelease(_ string) (string, error) {
	var res syscall.Utsname
	if err := syscall.Uname(&res); err != nil {
		return "", fmt.Errorf("unable to get uname: %w", err)
	}
	bytes := make([]byte, len(res.Release))
	for i, chr := range res.Release {
		bytes[i] = byte(chr)
	}
	return string(bytes), nil
}
