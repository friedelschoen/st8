package component

import (
	"fmt"
	"os"
)

func NumFiles(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("unable to read directory: %w", err)
	}
	return fmt.Sprintf("%d", len(entries)), nil
}
