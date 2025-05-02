package component

import (
	"os"

	"github.com/friedelschoen/st8/notify"
)

func Hostname(_ string, _ *notify.Notification, _ *any) (string, error) {
	return os.Hostname()
}
