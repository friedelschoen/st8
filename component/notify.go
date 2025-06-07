package component

import (
	"fmt"
	"strings"

	"github.com/friedelschoen/st8/notify"
)

func notifyAppName(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		if not == nil {
			return fmt.Errorf("cannot use notify_* in regular status")
		}
		block.Text = not.AppName
		return nil
	}, nil
}

func notifyAppIcon(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		if not == nil {
			return fmt.Errorf("cannot use notify_* in regular status")
		}
		block.Text = not.AppIcon
		return nil
	}, nil
}

func notifySummary(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		if not == nil {
			return fmt.Errorf("cannot use notify_* in regular status")
		}
		block.Text = not.Summary
		return nil
	}, nil
}

func notifyBody(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		if not == nil {
			return fmt.Errorf("cannot use notify_* in regular status")
		}
		block.Text = not.Body
		return nil
	}, nil
}

func notifyActions(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		if not == nil {
			return fmt.Errorf("cannot use notify_* in regular status")
		}
		block.Text = strings.Join(not.Actions, ", ")
		return nil
	}, nil
}

func init() {
	Install("notify_appname", notifyAppName)
	Install("notify_appicon", notifyAppIcon)
	Install("notify_summary", notifySummary)
	Install("notify_body", notifyBody)
	Install("notify_actions", notifyActions)
}
