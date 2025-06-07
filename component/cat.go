package component

import (
	"os"
	"strings"

	"github.com/friedelschoen/st8/notify"
)

func readFile(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		content, err := os.ReadFile(args["file"])
		block.Text = strings.TrimSpace(string(content))
		return err
	}, nil
}

func init() {
	Install("cat", readFile)
}
