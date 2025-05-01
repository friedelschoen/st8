package component

import (
	"fmt"

	"github.com/mdlayher/wifi"
)

func WifiESSID(interfaceName string) (string, error) {
	client, err := wifi.New()
	if err != nil {
		return "", fmt.Errorf("failed to open netlink: %w", err)
	}
	defer client.Close()

	interfaces, err := client.Interfaces()
	if err != nil {
		return "", fmt.Errorf("failed to list interfaces: %w", err)
	}

	var ifi *wifi.Interface
	for _, i := range interfaces {
		if i.Name == interfaceName {
			ifi = i
			break
		}
	}
	if ifi == nil {
		return "", fmt.Errorf("interface %q not found", interfaceName)
	}

	bss, err := client.BSS(ifi)
	if err != nil {
		return "", fmt.Errorf("failed to get BSS: %w", err)
	}
	return string(bss.SSID), nil
}

func WifiPerc(interfaceName string) (string, error) {
	client, err := wifi.New()
	if err != nil {
		return "", fmt.Errorf("failed to open netlink: %w", err)
	}
	defer client.Close()

	interfaces, err := client.Interfaces()
	if err != nil {
		return "", fmt.Errorf("failed to list interfaces: %w", err)
	}

	var ifi *wifi.Interface
	for _, i := range interfaces {
		if i.Name == interfaceName {
			ifi = i
			break
		}
	}
	if ifi == nil {
		return "", fmt.Errorf("interface %q not found", interfaceName)
	}

	stations, err := client.StationInfo(ifi)
	if err != nil || len(stations) == 0 {
		return "", fmt.Errorf("failed to get station info: %w", err)
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

	return fmt.Sprintf("%d", perc), nil
}
