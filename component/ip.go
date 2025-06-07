package component

import (
	"fmt"
	"net"
	"strings"

	"github.com/friedelschoen/st8/notify"
)

func getIPAddrs(ifaceName string) (ipv4s []string, ipv6s []string, err error) {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, nil, fmt.Errorf("interface not found: %w", err)
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get addresses: %w", err)
	}

	for _, addr := range addrs {
		var ip net.IP

		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}

		if ip == nil || ip.IsLoopback() {
			continue
		}
		if ip.To4() != nil {
			ipv4s = append(ipv4s, ip.String())
		} else if ip.To16() != nil {
			ipv6s = append(ipv6s, ip.String())
		}
	}
	return
}

func IPv4(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		ipv4s, _, err := getIPAddrs(args["interface"])
		if err != nil {
			return err
		}
		block.Text = strings.Join(ipv4s, ", ")
		return nil
	}, nil
}

func IPv6(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		_, ipv6s, err := getIPAddrs(args["interface"])
		if err != nil {
			return err
		}
		block.Text = strings.Join(ipv6s, ", ")
		return nil
	}, nil
}

func Up(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		netIface, err := net.InterfaceByName(args["interface"])
		if err != nil {
			return fmt.Errorf("interface not found: %w", err)
		}
		if netIface.Flags&net.FlagUp != 0 {
			block.Text = "up"
			return nil
		}
		block.Text = "down"
		return nil
	}, nil
}
