package component

import (
	"sync"
	"time"

	"github.com/friedelschoen/st8/notify"
	"github.com/shirou/gopsutil/v3/net"
)

type netstat struct {
	recv uint64
	time time.Time
}

var (
	lastRx = map[string]netstat{}
	lastTx = map[string]netstat{}
	mu     sync.Mutex
)

func NetspeedRx(interfaceName string, _ *notify.Notification) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	stats, err := net.IOCounters(true)
	if err != nil {
		return "", err
	}

	var rx uint64
	for _, s := range stats {
		if s.Name == interfaceName {
			rx = s.BytesRecv
			break
		}
	}
	now := time.Now()
	old, ok := lastRx[interfaceName]
	lastRx[interfaceName] = netstat{rx, now}
	if !ok || old.recv == 0 || now.Sub(old.time).Milliseconds() == 0 {
		return "0 B/s", nil // skip first read
	}
	diff := rx - old.recv
	bps := uint64(time.Second * time.Duration(diff) / now.Sub(old.time))
	return fmtHuman(bps) + "/s", nil
}

func NetspeedTx(interfaceName string, _ *notify.Notification) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	stats, err := net.IOCounters(true)
	if err != nil {
		return "", err
	}

	var tx uint64
	for _, s := range stats {
		if s.Name == interfaceName {
			tx = s.BytesSent
			break
		}
	}
	now := time.Now()
	old, ok := lastTx[interfaceName]
	lastTx[interfaceName] = netstat{tx, now}
	if !ok || old.recv == 0 || now.Sub(old.time).Milliseconds() == 0 {
		return "0 B/s", nil // skip first read
	}
	diff := tx - old.recv
	bps := uint64(time.Second * time.Duration(diff) / now.Sub(old.time))
	return fmtHuman(bps) + "/s", nil
}
