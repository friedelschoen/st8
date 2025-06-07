package component

import (
	"fmt"
	"os"
	"os/user"

	"github.com/friedelschoen/st8/notify"
)

func Gid(args map[string]string, events *EventHandlers) (Component, error) {
	gidtext := fmt.Sprintf("%d", os.Getgid())
	return func(block *Block, not *notify.Notification) error {
		block.Text = gidtext
		return nil
	}, nil
}

func Uid(args map[string]string, events *EventHandlers) (Component, error) {
	uidtext := fmt.Sprintf("%d", os.Getuid())
	return func(block *Block, not *notify.Notification) error {
		block.Text = uidtext
		return nil
	}, nil
}

func Username(args map[string]string, events *EventHandlers) (Component, error) {
	user, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("unable to determine user: %w", err)
	}
	return func(block *Block, not *notify.Notification) error {
		block.Text = user.Username
		return nil
	}, nil
}

func init() {
	Install("gid", Gid)
	Install("uid", Uid)
	Install("username", Username)
}
