package component

import (
	"log"
	"strings"

	"github.com/friedelschoen/ctxmenu"
	"github.com/friedelschoen/st8/notify"
	"github.com/friedelschoen/st8/proto"
	"github.com/friedelschoen/st8/sni"
	"github.com/godbus/dbus/v5"
)

var (
	conn *dbus.Conn
	host sni.TrayHost
)

func systray(args map[string]string, events *proto.EventHandlers) (Component, error) {
	if conn == nil {
		var err error
		conn, err = dbus.ConnectSessionBus()
		if err != nil {
			return nil, err
		}

		if _, err = sni.StartWatcher(conn); err != nil {
			return nil, err
		}

		if err := host.Run(conn); err != nil {
			return nil, err
		}
	}
	if filter, _ := args["filter"]; filter != "" {
		return filteredSystray(args, events)
	} else {
		return unfilteredSystray(args, events)
	}
}

func filteredSystray(args map[string]string, events *proto.EventHandlers) (Component, error) {
	filter, _ := args["filter"]

	events.OnClick = func(evt proto.ClickEvent) {
		for _, item := range host.Items {
			id, err := item.Id()
			if err != nil {
				log.Println(err)
				return
			}
			if !strings.Contains(id, filter) {
				continue
			}
			err = item.ContextMenu(evt.X, evt.Y)
			if err != nil {
				log.Println(err)
				return
			}
			break
		}
	}

	return func(block *proto.Block, not *notify.Notification) error {
		for _, item := range host.Items {
			id, err := item.Id()
			if err != nil {
				return err
			}
			if !strings.Contains(id, filter) {
				continue
			}
			title, err := item.Id()
			if err != nil {
				return err
			}
			block.Text = title
			break
		}
		return nil
	}, nil
}

func unfilteredSystray(args map[string]string, events *proto.EventHandlers) (Component, error) {
	remove, _ := args["remove"]

	events.OnClick = func(evt proto.ClickEvent) {
		conf := ctxmenu.Config{
			/* font, separate different fonts with comma */
			FontName: "monospace:size=12",

			/* colors */
			BackgroundColor:    "#FFFFFF",
			ForegroundColor:    "#2E3436",
			SelbackgroundColor: "#3584E4",
			SelforegroundColor: "#FFFFFF",
			SeparatorColor:     "#CDC7C2",
			BorderColor:        "#E6E6E6",

			/* sizes in pixels */
			MinItemWidth:    130, /* minimum width of a menu */
			BorderSize:      1,   /* menu border */
			SeperatorLength: 3,   /* space around separator */

			/* text alignment, set to LeftAlignment, CenterAlignment or RightAlignment */
			Alignment: ctxmenu.AlignLeft,

			/*
			 * The variables below cannot be set by X resources.
			 * Their values must be less than .height_pixels.
			 */

			/* the icon size is equal to .height_pixels - .iconpadding * 2 */
			IconSize: 24,

			/* area around the icon, the triangle and the separator */
			PaddingX: 4,
			PaddingY: 4,
		}

		menu := ctxmenu.Menu[int]{}
		for i, item := range host.Items {
			if remove != "" {
				id, err := item.Id()
				if err != nil {
					log.Printf("error: %v", err)
					return
				}
				doskip := false
				for f := range strings.FieldsSeq(remove) {
					if strings.Contains(id, f) {
						doskip = true
						break
					}
				}
				if doskip {
					continue
				}
			}
			title, err := item.Id()
			if err != nil {
				log.Println(err)
				return
			}
			menu = append(menu, ctxmenu.Item[int]{
				Label:  title,
				Output: i,
			})
		}

		id, err := ctxmenu.Run(menu, conf, "", nil)
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
		first := true
		for _, item := range host.Items {
			if remove != "" {
				id, err := item.Id()
				if err != nil {
					return err
				}
				doskip := false
				for f := range strings.FieldsSeq(remove) {
					if strings.Contains(id, f) {
						doskip = true
						break
					}
				}
				if doskip {
					continue
				}
			}
			if !first {
				block.Text += ", "
			}
			first = false
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
