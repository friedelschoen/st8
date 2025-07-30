package sni

import (
	"slices"
	"strings"

	"github.com/godbus/dbus/v5"
)

func getProp[T any](obj dbus.BusObject, iface string, prop string) (res T, err error) {
	var variant dbus.Variant
	if err = obj.Call("org.freedesktop.DBus.Properties.Get", 0, iface, prop).Store(&variant); err != nil {
		return res, err
	}
	res = variant.Value().(T)
	return
}

type TrayHost struct {
	Items []*TrayItem
}

func (host *TrayHost) Run(conn *dbus.Conn) error {
	watcher := conn.Object(watcherName, dbus.ObjectPath(watcherPath))
	busName := conn.Names()[0]

	err := watcher.Call(watcherInterface+".RegisterStatusNotifierHost", 0, busName).Err
	if err != nil {
		return err
	}

	ch := make(chan *dbus.Signal, 10)
	conn.Signal(ch)

	conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
		"type='signal',interface='org.kde.StatusNotifierWatcher',member='StatusNotifierItemRegistered'")
	conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
		"type='signal',interface='org.kde.StatusNotifierWatcher',member='StatusNotifierItemUnregistered'")

	variant, err := getProp[[]string](conn.BusObject(), watcherInterface, "RegisteredStatusNotifierItems")
	if err == nil {
		for _, service := range variant {
			host.Items = append(host.Items, newItem(conn, service))
		}
	}

	go func() {
		for sig := range ch {
			switch sig.Name {
			case watcherInterface + ".StatusNotifierItemRegistered":
				service := sig.Body[0].(string)
				if !slices.ContainsFunc(host.Items, func(i *TrayItem) bool { return i.service == service }) {
					host.Items = append(host.Items, newItem(conn, service))
				}
			case watcherInterface + ".StatusNotifierItemUnregistered":
				service := sig.Body[0].(string)
				host.Items = slices.DeleteFunc(host.Items, func(i *TrayItem) bool { return i.service == service })
			}
		}
	}()
	return nil
}

func newItem(conn *dbus.Conn, id string) *TrayItem {
	var item TrayItem
	item.service = id
	idx := strings.IndexByte(id, '/')
	item.obj = conn.Object(id[:idx], dbus.ObjectPath(id[idx:]))
	return &item
}

type TrayItem struct {
	service string
	obj     dbus.BusObject
}

func (item *TrayItem) Id() (string, error) {
	return getProp[string](item.obj, "org.kde.StatusNotifierItem", "Id")
}

func (item *TrayItem) Title() (string, error) {
	return getProp[string](item.obj, "org.kde.StatusNotifierItem", "Title")
}

func (item *TrayItem) Status() (string, error) {
	return getProp[string](item.obj, "org.kde.StatusNotifierItem", "Status")
}

func (item *TrayItem) ContextMenu(x, y int) error {
	res := item.obj.Call("org.kde.StatusNotifierItem.ContextMenu", 0, x, y)
	<-res.Done
	return res.Err
}
