package component

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/friedelschoen/st8/notify"
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

// rx-bytes: 1
// rx-packets: 2
// tx-bytes: 9
// tx-packets: 10
func getStat(pattern string, field int) (uint64, error) {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return 0, err
	}
	scanner := bufio.NewScanner(file)
	var total uint64
	for scanner.Scan() {
		line := scanner.Text()
		if strings.ContainsRune(line, '|') { // is header
			continue
		}
		fields := strings.Fields(line)
		name := strings.TrimSuffix(fields[0], ":")
		if !match(pattern, name) {
			continue
		}
		nr, err := strconv.Atoi(fields[field])
		if err != nil {
			return 0, err
		}
		total += uint64(nr)
	}
	return total, scanner.Err()
}

func NetspeedRx(block *Block, args map[string]string, not *notify.Notification, cacheptr *any) error {
	var cache netstat
	if *cacheptr != nil {
		cache = (*cacheptr).(netstat)
	}

	rx, err := getStat(args["interface"], 1)
	if err != nil {
		return err
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

	tx, err := getStat(args["interface"], 9)
	if err != nil {
		return err
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
