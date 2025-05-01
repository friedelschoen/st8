package component

import "os"

func Hostname(_ string) (string, error) {
	return os.Hostname()
}
