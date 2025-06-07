package component

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/friedelschoen/st8/notify"
)

func uptime(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		contents, err := os.ReadFile("/proc/uptime")
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
	}, nil
}

func init() {
	Install("uptime", uptime)
}
