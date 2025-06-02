package component

import (
	"fmt"
	"os"

	"github.com/friedelschoen/st8/notify"
)

func NumFiles(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	entries, err := os.ReadDir(args["path"])
	if err != nil {
		return fmt.Errorf("unable to read directory: %w", err)
	}
	block.Text = fmt.Sprintf("%d", len(entries))
	return nil
}
