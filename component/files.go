package component

import (
	"fmt"
	"os"

	"github.com/friedelschoen/st8/notify"
)

func NumFiles(dir string, _ *notify.Notification) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("unable to read directory: %w", err)
	}
	return fmt.Sprintf("%d", len(entries)), nil
}
