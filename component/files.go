package component

import (
	"fmt"
	"os"

	"github.com/friedelschoen/st8/notify"
	"github.com/friedelschoen/st8/proto"
)

func numFiles(args map[string]string, events *proto.EventHandlers) (Component, error) {
	path, ok := args["path"]
	if !ok {
		return nil, fmt.Errorf("missing argument: path")
	}
	return func(block *proto.Block, not *notify.Notification) error {
		entries, err := os.ReadDir(path)
		if err != nil {
			return fmt.Errorf("unable to read directory: %w", err)
		}
		block.Text = fmt.Sprintf("%d", len(entries))
		return nil
	}, nil
}

func init() {
	Install("num_files", numFiles)
}
