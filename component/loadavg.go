package component

import (
	"fmt"
	"os"
	"strings"

	"github.com/friedelschoen/st8/notify"
)

func LoadAverage(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
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
}

//number of processes currently runnable (running or on ready queue); total number of processes in system; last pid created. All fields are separated by one space except “number of processes currently runnable” and “total number of processes in system”, which are separated by a slash (‘/’). Example: 0.61 0.61 0.55 3/828 22084
