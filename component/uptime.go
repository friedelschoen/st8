package component

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/friedelschoen/st8/notify"
)

func Uptime(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	contents, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return err
	}
	fields := strings.Fields(string(contents))

	seconds, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return err
	}
	block.Text = time.Duration(seconds * float64(time.Second)).String()
	return nil
}
