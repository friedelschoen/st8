package popupmenu

import (
	"fmt"
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

	var menu ctxmenu.Menu[int]
	// Alleen de children, root zelf is meestal een dummy container
	rootitem, err := parseNode(root)
	if err != nil {
		return fmt.Errorf("unable to create menu: %w", err)
	}
	if rootitem.Label == "" {
		menu = rootitem.SubMenu
	}

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
	clicked, err := ctxmenu.Run(menu, conf, "", nil)
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

func parseNode(item dbusitem) (ctxmenu.Item[int], error) {
	var res ctxmenu.Item[int]
	if !getProp(item.Props, "visible", true) {
		return res, nil
	}
	itemType := getProp(item.Props, "type", "")
	if itemType == "separator" {
		return res, nil
	}

	res.Output = int(item.Id)
	res.Label = cleanLabel(getProp(item.Props, "label", ""))
	// enabled := getBool(item.props, "enabled", true)

	for _, child := range item.Children {
		m, err := parseNode(child)
		if err != nil {
			return res, err
		}
		res.SubMenu = append(res.SubMenu, m)
	}

	// // Disabled items worden niet opgenomen in xmenu
	// if !enabled && len(subItems) == 0 {
	// 	return nil, nil
	// }

	return res, nil
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
