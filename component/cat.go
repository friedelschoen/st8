package component

import (
	"os"
	"strings"

	"github.com/friedelschoen/st8/notify"
)

func ReadFile(file string, _ *notify.Notification) (string, error) {
	content, err := os.ReadFile(file)
	return strings.Trim(string(content), "\n\t\r "), err
}
