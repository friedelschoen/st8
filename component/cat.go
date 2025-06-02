package component

import (
	"os"
	"strings"

	"github.com/friedelschoen/st8/notify"
)

func ReadFile(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	content, err := os.ReadFile(args["file"])
	block.Text = strings.TrimSpace(string(content))
	return err
}
