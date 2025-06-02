package component

import (
	"os"

	"github.com/friedelschoen/st8/notify"
)

func Hostname(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	var err error
	block.Text, err = os.Hostname()
	return err
}
