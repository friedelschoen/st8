package component

import (
	"os"
	"strings"

	"github.com/friedelschoen/st8/notify"
)

func ReadFile(file string, _ *notify.Notification, _ *any) (string, error) {
	content, err := os.ReadFile(file)
	return strings.TrimSpace(string(content)), err
}
