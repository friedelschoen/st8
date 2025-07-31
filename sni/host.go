package sni

import (
	"log"
	"slices"
	"strings"

	"github.com/friedelschoen/st8/popupmenu"
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

	// conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
	// 	"type='signal',interface='org.kde.StatusNotifierWatcher',member='StatusNotifierItemUnregistered'")

	if err := conn.AddMatchSignal(dbus.WithMatchInterface("org.kde.StatusNotifierWatcher"), dbus.WithMatchMember("StatusNotifierItemRegistered")); err != nil {
		log.Printf("unable to setup watch for register: %v", err)
	}
	if err := conn.AddMatchSignal(dbus.WithMatchInterface("org.kde.StatusNotifierWatcher"), dbus.WithMatchMember("StatusNotifierItemUnregistered")); err != nil {
		log.Printf("unable to setup watch for register: %v", err)
	}

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
				if !slices.ContainsFunc(host.Items, func(i *TrayItem) bool { return i.id == service }) {
					host.Items = append(host.Items, newItem(conn, service))
				}
			case watcherInterface + ".StatusNotifierItemUnregistered":
				service := sig.Body[0].(string)
				host.Items = slices.DeleteFunc(host.Items, func(i *TrayItem) bool { return i.id == service })
			}
		}
	}()
	return nil
}

type TrayItem struct {
	id      string
	dest    string
	service dbus.ObjectPath
	conn    *dbus.Conn
	obj     dbus.BusObject
}

func newItem(conn *dbus.Conn, id string) *TrayItem {
	var item TrayItem
	item.id = id
	item.conn = conn
	idx := strings.IndexByte(id, '/')
	if idx == -1 {
		return nil
	}
	item.dest = id[:idx]
	item.service = dbus.ObjectPath(id[idx:])
	item.obj = item.conn.Object(item.dest, item.service)
	return &item
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
	call := item.obj.Call("org.kde.StatusNotifierItem.ContextMenu", 0, x, y)
	if d, ok := call.Err.(dbus.Error); !ok || d.Name != "org.freedesktop.DBus.Error.UnknownMethod" {
		return call.Err // <- unknown/other error
	}

	// busctl --user call :1.54 /org/blueman/sni/menu com.canonical.dbusmenu GetLayout iias -- 0 -1 0

	menupath, err := getProp[dbus.ObjectPath](item.obj, "org.kde.StatusNotifierItem", "Menu")
	if err != nil {
		return err
	}
	if err := popupmenu.NewDBusMenu(item.conn, item.dest, menupath); err != nil {
		return err
	}

	return nil
}
