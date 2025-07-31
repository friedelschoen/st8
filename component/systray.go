package component

import (
	"log"
	"strconv"

	"github.com/friedelschoen/st8/notify"
	"github.com/friedelschoen/st8/popupmenu"
	"github.com/friedelschoen/st8/proto"
	"github.com/friedelschoen/st8/sni"
	"github.com/godbus/dbus/v5"
)

func systray(args map[string]string, events *proto.EventHandlers) (Component, error) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return nil, err
	}

	if _, err = sni.StartWatcher(conn); err != nil {
		return nil, err
	}

	var host sni.TrayHost
	if err := host.Run(conn); err != nil {
		return nil, err
	}

	events.OnClick = func(evt proto.ClickEvent) {
		var items []popupmenu.MenuItem
		for i, item := range host.Items {
			title, err := item.Id()
			if err != nil {
				log.Println(err)
				return
			}
			items = append(items, popupmenu.MenuItem{Text: title, Id: strconv.Itoa(i)})
		}
		idstr, err := popupmenu.PopupMenu(items)
		if err != nil {
			log.Println(err)
			return
		}
		id, err := strconv.Atoi(idstr)
		if err != nil {
			log.Println(err)
			return
		}
		err = host.Items[id].ContextMenu(evt.X, evt.Y)
		if err != nil {
			log.Println(err)
			return
		}
	}

	return func(block *proto.Block, not *notify.Notification) error {
		for i, item := range host.Items {
			if i > 0 {
				block.Text += ", "
			}
			title, err := item.Id()
			if err != nil {
				return err
			}
			block.Text += title
		}
		return nil
	}, nil
}

func init() {
	Install("systray", systray)
}
