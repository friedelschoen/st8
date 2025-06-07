package component

import (
	"fmt"

	"github.com/friedelschoen/st8/notify"
	"github.com/mdlayher/wifi"
)

func WifiESSID(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		client, err := wifi.New()
		if err != nil {
			return fmt.Errorf("failed to open netlink: %w", err)
		}
		defer client.Close()

		interfaces, err := client.Interfaces()
		if err != nil {
			return fmt.Errorf("failed to list interfaces: %w", err)
		}

		var ifi *wifi.Interface
		for _, i := range interfaces {
			if globMatch(args["interface"], i.Name) {
				ifi = i
				break
			}
		}
		if ifi == nil {
			return fmt.Errorf("interface %q not found", args["interface"])
		}

		bss, err := client.BSS(ifi)
		if err != nil {
			return fmt.Errorf("failed to get BSS: %w", err)
		}
		block.Text = string(bss.SSID)
		return nil
	}, nil
}

func WifiPerc(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		client, err := wifi.New()
		if err != nil {
			return fmt.Errorf("failed to open netlink: %w", err)
		}
		defer client.Close()

		interfaces, err := client.Interfaces()
		if err != nil {
			return fmt.Errorf("failed to list interfaces: %w", err)
		}

		var ifi *wifi.Interface
		for _, i := range interfaces {
			if globMatch(args["interface"], i.Name) {
				ifi = i
				break
			}
		}
		if ifi == nil {
			return fmt.Errorf("interface %q not found", args["interface"])
		}

		stations, err := client.StationInfo(ifi)
		if err != nil || len(stations) == 0 {
			return fmt.Errorf("failed to get station info: %w", err)
		}

		rssi := stations[0].Signal
		// RSSI to percent (same as C macro)
		var perc int
		switch {
		case rssi >= -50:
			perc = 100
		case rssi <= -100:
			perc = 0
		default:
			perc = 2 * (rssi + 100)
		}

		block.Text = fmt.Sprintf("%d", perc)
		return nil
	}, nil
}

func init() {
	Install("wifi_essid", WifiESSID)
	Install("wifi_perc", WifiPerc)
}
