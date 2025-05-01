package main

import (
	"fmt"

	"github.com/godbus/dbus/v5"
)

var (
	notificationID uint32 = 1
)

type Notification struct {
	appName    string
	replacesID uint32
	appIcon    string
	summary    string
	body       string
	actions    []string
	hints      map[string]dbus.Variant
	timeout    int32
}

type NotificationDaemon struct {
	*dbus.Conn
	C chan Notification
}

// D-Bus: GetCapabilities method
func (n *NotificationDaemon) GetCapabilities() ([]string, *dbus.Error) {
	return []string{"body", "actions", "icon-static"}, nil
}

// D-Bus: GetServerInformation method
func (n *NotificationDaemon) GetServerInformation() (name, vendor, version, specVersion string, err *dbus.Error) {
	return "dwm-status", "friedelschoen", "0.1", "1.2", nil
}

// D-Bus: Notify method
func (n *NotificationDaemon) Notify(appName string, replacesID uint32, appIcon string, summary string, body string, actions []string, hints map[string]dbus.Variant, timeout int32) (uint32, *dbus.Error) {
	n.C <- Notification{
		appName, replacesID, appIcon, summary, body, actions, hints, timeout,
	}
	if replacesID == 0 {
		notificationID++
		return notificationID - 1, nil
	}
	return replacesID, nil
}

// D-Bus: CloseNotification method
func (n *NotificationDaemon) CloseNotification(id uint32) *dbus.Error {
	// hier kan je nog echte afsluitlogica implementeren
	fmt.Printf("CloseNotification called for id: %d\n", id)
	return nil
}

func NotifyStart() (*NotificationDaemon, error) {
	var conn NotificationDaemon

	conn.C = make(chan Notification)

	var err error
	conn.Conn, err = dbus.ConnectSessionBus()
	if err != nil {
		return nil, err
	}

	err = conn.Export(&conn, "/org/freedesktop/Notifications", "org.freedesktop.Notifications")
	if err != nil {
		return nil, err
	}

	reply, err := conn.RequestName("org.freedesktop.Notifications", dbus.NameFlagDoNotQueue)
	if err != nil {
		return nil, err
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		return nil, fmt.Errorf("another daemon running")
	}

	fmt.Println("Notification daemon is running...")

	return &conn, nil
}
