package component

import (
	"fmt"
	"os"
	"os/user"

	"github.com/friedelschoen/st8/notify"
	"github.com/friedelschoen/st8/proto"
)

func getGID(args map[string]string, events *proto.EventHandlers) (Component, error) {
	gidtext := fmt.Sprintf("%d", os.Getgid())
	return func(block *proto.Block, not *notify.Notification) error {
		block.Text = gidtext
		return nil
	}, nil
}

func getUID(args map[string]string, events *proto.EventHandlers) (Component, error) {
	uidtext := fmt.Sprintf("%d", os.Getuid())
	return func(block *proto.Block, not *notify.Notification) error {
		block.Text = uidtext
		return nil
	}, nil
}

func username(args map[string]string, events *proto.EventHandlers) (Component, error) {
	user, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("unable to determine user: %w", err)
	}
	return func(block *proto.Block, not *notify.Notification) error {
		block.Text = user.Username
		return nil
	}, nil
}

func init() {
	Install("gid", getGID)
	Install("uid", getUID)
	Install("username", username)
}
