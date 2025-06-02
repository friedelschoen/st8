package component

import (
	"fmt"
	"os"
	"os/user"

	"github.com/friedelschoen/st8/notify"
)

func Gid(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	gid := os.Getgid()
	block.Text = fmt.Sprintf("%d", gid)
	return nil
}

func Uid(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	uid := os.Getuid()
	block.Text = fmt.Sprintf("%d", uid)
	return nil
}

func Username(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	user, err := user.Current()
	if err != nil {
		return fmt.Errorf("unable to determine user: %w", err)
	}
	block.Text = user.Username
	return nil
}
