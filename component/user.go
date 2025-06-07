package component

import (
	"fmt"
	"os"
	"os/user"

	"github.com/friedelschoen/st8/notify"
)

func Gid(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		gid := os.Getgid()
		block.Text = fmt.Sprintf("%d", gid)
		return nil
	}, nil
}

func Uid(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		uid := os.Getuid()
		block.Text = fmt.Sprintf("%d", uid)
		return nil
	}, nil
}

func Username(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		user, err := user.Current()
		if err != nil {
			return fmt.Errorf("unable to determine user: %w", err)
		}
		block.Text = user.Username
		return nil
	}, nil
}

func init() {
	Install("gid", Gid)
	Install("uid", Uid)
	Install("username", Username)
}
