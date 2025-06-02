package component

import (
	"fmt"
	"time"

	"github.com/friedelschoen/st8/notify"
	"github.com/lestrrat-go/strftime"
)

var dateformats = make(map[string]*strftime.Strftime)

func Datetime(block *Block, args map[string]string, not *notify.Notification, cacheptr *any) error {
	var datefmt *strftime.Strftime
	if *cacheptr != nil {
		datefmt = (*cacheptr).(*strftime.Strftime)
	} else {
		var err error
		datefmt, err = strftime.New(args["datefmt"])
		if err != nil {
			return fmt.Errorf("unable to parse date-format `%s`: %w", args["datefmt"], err)
		}
		*cacheptr = datefmt
	}

	block.Text = datefmt.FormatString(time.Now())
	return nil
}
