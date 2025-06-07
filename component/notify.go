package component

import (
	"fmt"
	"strings"

	"github.com/friedelschoen/st8/notify"
)

func NotifyAppName(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		if not == nil {
			return fmt.Errorf("cannot use notify_* in regular status")
		}
		block.Text = not.AppName
		return nil
	}, nil
}

func NotifyAppIcon(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		if not == nil {
			return fmt.Errorf("cannot use notify_* in regular status")
		}
		block.Text = not.AppIcon
		return nil
	}, nil
}

func NotifySummary(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		if not == nil {
			return fmt.Errorf("cannot use notify_* in regular status")
		}
		block.Text = not.Summary
		return nil
	}, nil
}

func NotifyBody(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		if not == nil {
			return fmt.Errorf("cannot use notify_* in regular status")
		}
		block.Text = not.Body
		return nil
	}, nil
}

func NotifyActions(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		if not == nil {
			return fmt.Errorf("cannot use notify_* in regular status")
		}
		block.Text = strings.Join(not.Actions, ", ")
		return nil
	}, nil
}
