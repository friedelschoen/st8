package component

import (
	"fmt"
	"os"
	"strings"

	"github.com/friedelschoen/st8/notify"
	"github.com/friedelschoen/st8/proto"
)

func loadAverage(args map[string]string, events *proto.EventHandlers) (Component, error) {
	return func(block *proto.Block, not *notify.Notification) error {
		contents, err := os.ReadFile("/proc/loadavg")
		if err != nil {
			return err
		}
		fields := strings.Fields(string(contents))

		switch args["minutes"] {
		case "":
			block.Text = fmt.Sprintf("%s %s %s", fields[0], fields[1], fields[2])
			return nil
		case "1":
			block.Text = fields[0]
			return nil
		case "5":
			block.Text = fields[1]
			return nil
		case "15":
			block.Text = fields[2]
			return nil
		default:
			return fmt.Errorf("unable to get average load: period must be either 1, 5 or 15 (minutes)")
		}
	}, nil
}

func init() {
	Install("load_avg", loadAverage)
}
