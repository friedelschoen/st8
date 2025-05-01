package component

import (
	"fmt"
	"os"
	"os/user"
)

func Gid(_ string) (string, error) {
	gid := os.Getgid()
	return fmt.Sprintf("%d", gid), nil
}

func Uid(_ string) (string, error) {
	uid := os.Getuid()
	return fmt.Sprintf("%d", uid), nil
}

func Username(_ string) (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("unable to determine user: %w", err)
	}
	return user.Username, nil
}
