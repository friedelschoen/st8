package popupmenu

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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

	var menu []MenuItem
	// Alleen de children, root zelf is meestal een dummy container
	rootitem, err := parseNode(root)
	if err != nil {
		return fmt.Errorf("unable to create menu: %w", err)
	}
	if rootitem != nil {
		menu = rootitem.Children
	}

	clicked, err := PopupMenu(menu)
	if err != nil {
		return fmt.Errorf("unable to open menu: %w", err)
	}
	if len(clicked) == 0 {
		return nil
	}

	clickID, err := strconv.Atoi(clicked)
	if err != nil {
		return fmt.Errorf("invalid id `%s`: %w", clicked, err)
	}

	stamp := time.Now().Second()
	call = obj.Call("com.canonical.dbusmenu.Event", 0, int32(clickID), "clicked", dbus.MakeVariant(0), uint32(stamp))
	if call.Err != nil {
		return fmt.Errorf("Event failed: %w", call.Err)
	}
	return nil
}

func parseNode(item dbusitem) (*MenuItem, error) {
	if !getProp(item.Props, "visible", true) {
		return nil, nil
	}
	itemType := getProp(item.Props, "type", "")
	if itemType == "separator" {
		return &MenuItem{Text: "―", Id: ""}, nil
	}

	var res MenuItem
	res.Id = strconv.Itoa(int(item.Id))
	res.Text = cleanLabel(getProp(item.Props, "label", ""))
	// enabled := getBool(item.props, "enabled", true)

	for _, child := range item.Children {
		m, err := parseNode(child)
		if err != nil {
			return nil, err
		}
		if m != nil {
			res.Children = append(res.Children, *m)
		}
	}

	// // Disabled items worden niet opgenomen in xmenu
	// if !enabled && len(subItems) == 0 {
	// 	return nil, nil
	// }

	return &res, nil
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
