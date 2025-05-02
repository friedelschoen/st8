package component

import (
	"time"

	"github.com/friedelschoen/st8/notify"
	"github.com/shirou/gopsutil/v3/net"
)

type netstat struct {
	recv uint64
	time time.Time
}

func NetspeedRx(interfaceName string, _ *notify.Notification, cacheptr *any) (string, error) {
	var cache netstat
	if *cacheptr != nil {
		cache = (*cacheptr).(netstat)
	}

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
	*cacheptr = netstat{rx, now}
	if cache.recv == 0 || now.Sub(cache.time).Milliseconds() == 0 {
		return "0 B/s", nil // skip first read
	}
	diff := rx - cache.recv
	bps := uint64(time.Second * time.Duration(diff) / now.Sub(cache.time))
	return fmtHuman(bps) + "/s", nil
}

func NetspeedTx(interfaceName string, _ *notify.Notification, cacheptr *any) (string, error) {
	var cache netstat
	if *cacheptr != nil {
		cache = (*cacheptr).(netstat)
	}

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
	*cacheptr = netstat{tx, now}
	if cache.recv == 0 || now.Sub(cache.time).Milliseconds() == 0 {
		return "0 B/s", nil // skip first read
	}
	diff := tx - cache.recv
	bps := uint64(time.Second * time.Duration(diff) / now.Sub(cache.time))
	return fmtHuman(bps) + "/s", nil
}
