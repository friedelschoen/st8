package component

import (
	"fmt"
	"strings"

	"github.com/friedelschoen/st8/notify"
)

func NotifyAppName(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	if not == nil {
		return fmt.Errorf("cannot use notify_* in regular status")
	}
	block.Text = not.AppName
	return nil
}

func NotifyAppIcon(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	if not == nil {
		return fmt.Errorf("cannot use notify_* in regular status")
	}
	block.Text = not.AppIcon
	return nil
}

func NotifySummary(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	if not == nil {
		return fmt.Errorf("cannot use notify_* in regular status")
	}
	block.Text = not.Summary
	return nil
}

func NotifyBody(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	if not == nil {
		return fmt.Errorf("cannot use notify_* in regular status")
	}
	block.Text = not.Body
	return nil
}

func NotifyActions(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	if not == nil {
		return fmt.Errorf("cannot use notify_* in regular status")
	}
	block.Text = strings.Join(not.Actions, ", ")
	return nil
}
