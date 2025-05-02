package component

import (
	"fmt"
	"strings"

	"github.com/friedelschoen/st8/notify"
)

func NotifyAppName(_ string, not *notify.Notification, _ *any) (string, error) {
	if not == nil {
		return "", fmt.Errorf("cannot use notify_* in regular status")
	}
	return not.AppName, nil
}

func NotifyAppIcon(_ string, not *notify.Notification, _ *any) (string, error) {
	if not == nil {
		return "", fmt.Errorf("cannot use notify_* in regular status")
	}
	return not.AppIcon, nil
}

func NotifySummary(_ string, not *notify.Notification, _ *any) (string, error) {
	if not == nil {
		return "", fmt.Errorf("cannot use notify_* in regular status")
	}
	return not.Summary, nil
}

func NotifyBody(_ string, not *notify.Notification, _ *any) (string, error) {
	if not == nil {
		return "", fmt.Errorf("cannot use notify_* in regular status")
	}
	return not.Body, nil
}

func NotifyActions(_ string, not *notify.Notification, _ *any) (string, error) {
	if not == nil {
		return "", fmt.Errorf("cannot use notify_* in regular status")
	}
	return strings.Join(not.Actions, ", "), nil
}
