package component

import (
	"time"

	"github.com/friedelschoen/st8/notify"
	"github.com/ncruces/go-strftime"
)

func Datetime(block *Block, args map[string]string, not *notify.Notification, cacheptr *any) error {
	block.Text = strftime.Format(args["datefmt"], time.Now())
	return nil
}
