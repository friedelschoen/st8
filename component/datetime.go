package component

import (
	"fmt"
	"time"

	"github.com/friedelschoen/st8/notify"
	"github.com/lestrrat-go/strftime"
)

var dateformats = make(map[string]*strftime.Strftime)

func Datetime(format string, _ *notify.Notification) (string, error) {
	datefmt, ok := dateformats[format]
	if !ok {
		var err error
		datefmt, err = strftime.New(format)
		if err != nil {
			return "", fmt.Errorf("unable to parse date-format `%s`: %w", format, err)
		}
		dateformats[format] = datefmt
	}
	return datefmt.FormatString(time.Now()), nil
}
