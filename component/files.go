package component

import (
	"fmt"
	"os"

	"github.com/friedelschoen/st8/notify"
)

func NumFiles(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		entries, err := os.ReadDir(args["path"])
		if err != nil {
			return fmt.Errorf("unable to read directory: %w", err)
		}
		block.Text = fmt.Sprintf("%d", len(entries))
		return nil
	}, nil
}
