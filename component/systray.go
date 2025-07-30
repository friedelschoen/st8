package component

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/friedelschoen/st8/notify"
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
		cmd := exec.Command("xmenu")
		in, err := cmd.StdinPipe()
		var out strings.Builder
		cmd.Stdout = &out
		if err != nil {
			log.Println(err)
			return
		}
		if err := cmd.Start(); err != nil {
			log.Println(err)
			return
		}
		for i, item := range host.Items {
			title, err := item.Id()
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Fprintf(in, "%s\t%d\n", title, i)
		}
		in.Close()
		if err := cmd.Wait(); err != nil {
			log.Println(err)
			return
		}
		index, err := strconv.Atoi(strings.TrimSpace(out.String()))
		if err != nil {
			log.Println(err)
			return
		}
		host.Items[index].ContextMenu(evt.X, evt.Y)
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
