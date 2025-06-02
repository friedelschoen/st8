package component

import (
	"strings"
	"time"

	"github.com/friedelschoen/st8/notify"
	"github.com/shirou/gopsutil/v3/net"
)

type netstat struct {
	recv uint64
	time time.Time
}

func match(pattern, text string) bool {
	prefix, suffix, ok := strings.Cut(pattern, "*")
	if ok {
		return strings.HasPrefix(text, prefix) && strings.HasSuffix(text, suffix)
	}
	return pattern == text
}

func NetspeedRx(block *Block, args map[string]string, not *notify.Notification, cacheptr *any) error {
	var cache netstat
	if *cacheptr != nil {
		cache = (*cacheptr).(netstat)
	}

	stats, err := net.IOCounters(true)
	if err != nil {
		return err
	}

	var rx uint64
	for _, s := range stats {
		if match(args["interface"], s.Name) {
			rx += s.BytesRecv
		}
	}
	now := time.Now()
	*cacheptr = netstat{rx, now}
	if cache.recv == 0 || now.Sub(cache.time).Milliseconds() == 0 {
		block.Text = "0 B/s"
		return nil // skip first read
	}
	diff := rx - cache.recv
	bps := uint64(time.Second * time.Duration(diff) / now.Sub(cache.time))
	block.Text = fmtHuman(bps) + "/s"
	return nil
}

func NetspeedTx(block *Block, args map[string]string, not *notify.Notification, cacheptr *any) error {
	var cache netstat
	if *cacheptr != nil {
		cache = (*cacheptr).(netstat)
	}

	stats, err := net.IOCounters(true)
	if err != nil {
		return err
	}

	var tx uint64
	for _, s := range stats {
		if match(args["interface"], s.Name) {
			tx += s.BytesSent
		}
	}
	now := time.Now()
	*cacheptr = netstat{tx, now}
	if cache.recv == 0 || now.Sub(cache.time).Milliseconds() == 0 {
		block.Text = "0 B/s"
		return nil // skip first read
	}
	diff := tx - cache.recv
	bps := uint64(time.Second * time.Duration(diff) / now.Sub(cache.time))
	block.Text = fmtHuman(bps) + "/s"
	return nil
}
