package popupmenu

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/friedelschoen/ctxmenu"
	"github.com/godbus/dbus/v5"
)

type dbusitem struct {
	Id       int32
	Props    map[string]dbus.Variant
	Children []dbusitem
}

func NewDBusMenu(conn *dbus.Conn, dest string, path dbus.ObjectPath) error {
	obj := conn.Object(dest, path)

	var revision uint32

	call := obj.Call("com.canonical.dbusmenu.GetLayout", 0,
		int32(0), int32(-1), []string{})
	if call.Err != nil {
		return fmt.Errorf("GetLayout failed: %w", call.Err)
	}
	var root dbusitem
	if err := call.Store(&revision, &root); err != nil {
		return fmt.Errorf("Store layout failed: %w", err)
	}

	// Alleen de children, root zelf is meestal een dummy container
	menu, err := parseNode([]dbusitem{root})
	if err != nil {
		return fmt.Errorf("unable to create menu: %w", err)
	}

	clicked, err := ctxmenu.Run(menu, nil, "", nil)
	if err != nil {
		return fmt.Errorf("unable to open menu: %w", err)
	}

	stamp := time.Now().Second()
	call = obj.Call("com.canonical.dbusmenu.Event", 0, int32(clicked), "clicked", dbus.MakeVariant(0), uint32(stamp))
	if call.Err != nil {
		return fmt.Errorf("Event failed: %w", call.Err)
	}
	return nil
}

func getProp[T any](m map[string]dbus.Variant, key string, def T) (res T) {
	if v, ok := m[key]; ok {
		if s, ok := v.Value().(T); ok {
			return s
		}
	}
	return def
}

func cleanLabel(label string) string {
	label = strings.ReplaceAll(label, "__", "\x00") // tijdelijke placeholder
	label = strings.ReplaceAll(label, "_", "")
	label = strings.ReplaceAll(label, "\x00", "_")
	label = strings.ReplaceAll(label, "\n", " ")
	return strings.TrimSpace(label)
}

func parseNode(items []dbusitem) (res []ctxmenu.Item[int], _ error) {
	for _, item := range items {
		if !getProp(item.Props, "visible", true) {
			continue
		}
		itemType := getProp(item.Props, "type", "")
		if itemType == "separator" {
			res = append(res, &ctxmenu.SeparatorItem[int]{})
			continue
		}

		label := cleanLabel(getProp(item.Props, "label", ""))

		if label == "" {
			m, err := parseNode(item.Children)
			if err != nil {
				return nil, err
			}
			fmt.Fprintf(os.Stderr, "id=%d, label=%s, children=%d -> %d\n", item.Id, label, len(item.Children), len(m))
			if len(m) > 0 {
				res = append(res, m...)
			} else {
				res = append(res, &ctxmenu.SeparatorItem[int]{})
			}
		} else {
			i := &ctxmenu.LabelItem[int]{
				Output: int(item.Id),
				Text:   label,
			}
			// enabled := getBool(item.props, "enabled", true)

			m, err := parseNode(item.Children)
			if err != nil {
				return nil, err
			}
			i.SubMenu = m
			res = append(res, i)
		}
	}
	return res, nil
}
