package sni

import (
	"fmt"
	"slices"
	"sync"

	"github.com/godbus/dbus/v5"
)

const (
	watcherName      = "org.kde.StatusNotifierWatcher"
	watcherPath      = "/StatusNotifierWatcher"
	watcherInterface = "org.kde.StatusNotifierWatcher"
	hostName         = "org.friedelschoen.st8"
)

type watcher struct {
	mu    sync.Mutex
	items []string
	hosts []string
	conn  *dbus.Conn
}

func (w *watcher) RegisterStatusNotifierItem(serviceOrPath string, sender dbus.Sender) *dbus.Error {
	if len(serviceOrPath) == 0 {
		return nil
	}
	w.mu.Lock()
	defer w.mu.Unlock()

	var service, path string
	if serviceOrPath[0] == '/' {
		service = string(sender)
		path = serviceOrPath
	} else {
		service = serviceOrPath
		path = "/StatusNotifierItem"
	}
	id := service + path
	if slices.Contains(w.items, id) {
		return nil
	}

	w.items = append(w.items, id)
	w.conn.Emit(dbus.ObjectPath(watcherPath),
		watcherInterface+".StatusNotifierItemRegistered", id)

	w.watchItem(service, id)

	return nil
}

func (w *watcher) RegisterStatusNotifierHost(service string) *dbus.Error {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.hosts = append(w.hosts, service)
	w.watchItem(service, "")
	w.conn.Emit(dbus.ObjectPath(watcherPath),
		watcherInterface+".StatusNotifierHostRegistered")
	return nil
}

func (w *watcher) GetRegisteredStatusNotifierItems() ([]string, *dbus.Error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.items, nil
}

func (w *watcher) GetIsStatusNotifierHostRegistered() (bool, *dbus.Error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(w.hosts) != 0, nil
}

func (w *watcher) GetProtocolVersion() (int32, *dbus.Error) {
	return 0, nil
}

func (w *watcher) watchItem(service string, id string) {
	go func() {
		ch := make(chan *dbus.Signal, 1)
		w.conn.Signal(ch)

		match := fmt.Sprintf("type='signal',sender='org.freedesktop.DBus',interface='org.freedesktop.DBus',member='NameOwnerChanged',arg0='%s'", service)
		w.conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, match)

		for sig := range ch {
			if sig.Name == "org.freedesktop.DBus.NameOwnerChanged" {
				old := sig.Body[1].(string)
				new := sig.Body[2].(string)
				if old != "" && new == "" {
					w.mu.Lock()
					w.items = slices.DeleteFunc(w.items, func(s string) bool { return s == id })
					w.mu.Unlock()
					w.conn.Emit(dbus.ObjectPath(watcherPath),
						watcherInterface+".StatusNotifierItemUnregistered", id)
					return
				}
			}
		}
	}()
}

func StartWatcher(conn *dbus.Conn) (bool, error) {
	var watcher watcher
	watcher.conn = conn

	conn.Export(&watcher, dbus.ObjectPath(watcherPath), watcherInterface)

	reply, err := conn.RequestName(watcherName, dbus.NameFlagDoNotQueue)
	return reply == dbus.RequestNameReplyPrimaryOwner, err
}
