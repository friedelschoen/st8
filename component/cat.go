package component

import (
	"os"
	"strings"

	"github.com/friedelschoen/st8/notify"
	"github.com/friedelschoen/st8/proto"
)

func readFile(args map[string]string, events *proto.EventHandlers) (Component, error) {
	return func(block *proto.Block, not *notify.Notification) error {
		content, err := os.ReadFile(args["file"])
		block.Text = strings.TrimSpace(string(content))
		return err
	}, nil
}

func init() {
	Install("cat", readFile)
}
