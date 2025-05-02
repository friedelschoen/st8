package component

import (
	"time"

	"github.com/friedelschoen/st8/notify"
)

func Datetime(format string, _ *notify.Notification) (string, error) {
	return time.Now().Format(format), nil
}
