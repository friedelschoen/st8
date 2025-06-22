package notify

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/friedelschoen/st8/config"
	"github.com/godbus/dbus/v5"
)

var (
	currentID uint32 = 1
)

type DBusDaemon struct {
	*dbus.Conn
	C chan<- Notification
}

// D-Bus: GetCapabilities method
func (n *DBusDaemon) GetCapabilities() ([]string, *dbus.Error) {
	return []string{"body", "actions", "icon-static"}, nil
}

// D-Bus: GetServerInformation method
func (n *DBusDaemon) GetServerInformation() (name, vendor, version, specVersion string, err *dbus.Error) {
	return "st8", "friedelschoen", "0.1", "1.2", nil
}

// D-Bus: Notify method
func (n *DBusDaemon) Notify(appName string, replacesID uint32, appIcon string, summary string, body string, actions []string, _ map[string]dbus.Variant, timeout int32) (uint32, *dbus.Error) {
	var dur time.Duration
	if timeout > 0 {
		dur = time.Duration(timeout) * time.Millisecond
	}
	n.C <- Notification{
		appName, appIcon, summary, body, actions, dur,
	}
	if replacesID == 0 {
		currentID++
		return currentID - 1, nil
	}
	return replacesID, nil
}

// D-Bus: CloseNotification method
func (n *DBusDaemon) CloseNotification(id uint32) *dbus.Error {
	// hier kan je nog echte afsluitlogica implementeren
	fmt.Fprintf(os.Stderr, "CloseNotification called for id: %d\n", id)
	return nil
}

func startDBusDaemon(_ *config.MainConfig, channel chan<- Notification) (io.Closer, error) {
	var conn DBusDaemon

	conn.C = channel

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

	fmt.Fprintln(os.Stderr, "Notification daemon is running...")

	return &conn, nil
}

func init() {
	Install("dbus", startDBusDaemon)
}
