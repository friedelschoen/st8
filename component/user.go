package component

import (
	"fmt"
	"os"
	"os/user"

	"github.com/friedelschoen/st8/notify"
)

func Gid(_ string, _ *notify.Notification, _ *any) (string, error) {
	gid := os.Getgid()
	return fmt.Sprintf("%d", gid), nil
}

func Uid(_ string, _ *notify.Notification, _ *any) (string, error) {
	uid := os.Getuid()
	return fmt.Sprintf("%d", uid), nil
}

func Username(_ string, _ *notify.Notification, _ *any) (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("unable to determine user: %w", err)
	}
	return user.Username, nil
}
