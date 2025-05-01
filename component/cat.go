package component

import "os"

func ReadFile(file string) (string, error) {
	content, err := os.ReadFile(file)
	return string(content), err
}
