package main

import (
	"fmt"
	"log"
	"slices"
	"sync"
	"time"

	"github.com/godbus/dbus/v5"
)

const (
	watcherName      = "org.kde.StatusNotifierWatcher"
	watcherPath      = "/StatusNotifierWatcher"
	watcherInterface = "org.kde.StatusNotifierWatcher"
	hostName         = "org.friedelschoen.st8"
)

type Watcher struct {
	mu    sync.Mutex
	items []string
	hosts []string
	conn  *dbus.Conn
}

func NewWatcher(conn *dbus.Conn) *Watcher {
	return &Watcher{
		conn: conn,
	}
}

func (w *Watcher) RegisterStatusNotifierItem(path string) *dbus.Error {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.items = append(w.items, string(path))

	log.Printf("✅ Nieuw item: %s", string(path))

	w.conn.Emit(dbus.ObjectPath(watcherPath),
		watcherInterface+".StatusNotifierItemRegistered", string(path))

	w.watchItem(string(path))

	return nil
}

func (w *Watcher) RegisterStatusNotifierHost(service string) *dbus.Error {
	w.mu.Lock()
	defer w.mu.Unlock()

	log.Printf("✅ Host geregistreerd: %s", service)
	w.hosts = append(w.hosts, service)
	w.watchItem(service)
	w.conn.Emit(dbus.ObjectPath(watcherPath),
		watcherInterface+".StatusNotifierHostRegistered")
	return nil
}

func (w *Watcher) GetRegisteredStatusNotifierItems() ([]string, *dbus.Error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.items, nil
}

func (w *Watcher) GetIsStatusNotifierHostRegistered() (bool, *dbus.Error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(w.hosts) != 0, nil
}

func (w *Watcher) GetProtocolVersion() (int32, *dbus.Error) {
	return 0, nil
}

func (w *Watcher) watchItem(service string) {
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
					w.items = slices.DeleteFunc(w.items, func(s string) bool { return s == service })
					w.mu.Unlock()
					log.Printf("❌ Item verdween van bus: %s", service)
					w.conn.Emit(dbus.ObjectPath(watcherPath),
						watcherInterface+".StatusNotifierItemUnregistered", service)
					return
				}
			}
		}
	}()
}

func runHost(conn *dbus.Conn) {
	watcher := conn.Object(watcherName, dbus.ObjectPath(watcherPath))
	busName := conn.Names()[0]

	err := watcher.Call(watcherInterface+".RegisterStatusNotifierHost", 0, busName).Err
	if err != nil {
		log.Fatalf("Kon host niet registreren: %v", err)
	}
	fmt.Println("📡 Geregistreerd als StatusNotifierHost")

	// Signalen volgen
	ch := make(chan *dbus.Signal, 10)
	conn.Signal(ch)

	conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
		"type='signal',interface='org.kde.StatusNotifierWatcher',member='StatusNotifierItemRegistered'")

	// Initiale lijst ophalen
	var variant dbus.Variant
	err = watcher.Call("org.freedesktop.DBus.Properties.Get", 0,
		watcherInterface, "RegisteredStatusNotifierItems").Store(&variant)

	if err == nil {
		items := variant.Value().([]string)
		for _, s := range items {
			printItemInfo(conn, s)
		}
	}

	for sig := range ch {
		if sig.Name == watcherInterface+".StatusNotifierItemRegistered" {
			service := sig.Body[0].(string)
			fmt.Printf("➕ Nieuw item op de bus: %s\n", service)
			printItemInfo(conn, service)
		}
	}
}

func printItemInfo(conn *dbus.Conn, service string) {
	time.Sleep(100 * time.Millisecond) // wacht even tot item klaar is

	obj := conn.Object(service, dbus.ObjectPath("org.kde.StatusNotifierItem"))

	props := []string{"Title", "Id", "Status", "IconName"}

	fmt.Printf("📦 Item: %s\n", service)
	for _, p := range props {
		if val, err := getProp(obj, p); err == nil {
			fmt.Printf("  %s: %v\n", p, val)
		} else {
			fmt.Printf("  %s: ERR %v\n", p, err)
		}
	}
	fmt.Println()
}

func getProp(obj dbus.BusObject, prop string) (interface{}, error) {
	variant := dbus.Variant{}
	err := obj.Call("org.freedesktop.DBus.Properties.Get", 0, "org.kde.StatusNotifierItem", prop).Store(&variant)
	if err != nil {
		return nil, err
	}
	return variant.Value(), nil
}

func main() {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		log.Fatalf("dbus connect failed: %v", err)
	}
	defer conn.Close()

	// Export watcher
	watcher := NewWatcher(conn)
	conn.Export(watcher, dbus.ObjectPath(watcherPath), watcherInterface)

	reply, err := conn.RequestName(watcherName, dbus.NameFlagDoNotQueue)
	if err != nil {
		log.Fatalf("failed to request name: %v", err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		log.Printf("watcher already running")
	} else {
		log.Printf("watcher running")
	}

	runHost(conn)

	select {}
}
