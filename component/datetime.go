package component

import "time"

func Datetime(format string) (string, error) {
	return time.Now().Format(format), nil
}
