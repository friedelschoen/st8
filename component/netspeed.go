package component

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/friedelschoen/st8/notify"
)

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
		if !globMatch(pattern, name) {
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

func netspeedRx(args map[string]string, events *EventHandlers) (Component, error) {
	var recv uint64
	var lastTime time.Time

	return func(block *Block, not *notify.Notification) error {
		rx, err := getStat(args["interface"], 1)
		if err != nil {
			return err
		}

		now := time.Now()
		if recv == 0 {
			lastTime = now
			recv = rx
			block.Text = "0 B/s"
			return nil // skip first read
		}
		bps := float64(rx-recv) / now.Sub(lastTime).Seconds()
		lastTime = now
		recv = rx
		block.Text = fmtHuman(uint64(bps)) + "/s"
		return nil
	}, nil
}

func netspeedTx(args map[string]string, events *EventHandlers) (Component, error) {
	var sent uint64
	var lastTime time.Time

	return func(block *Block, not *notify.Notification) error {
		tx, err := getStat(args["interface"], 9)
		if err != nil {
			return err
		}

		now := time.Now()
		if sent == 0 {
			lastTime = now
			sent = tx
			block.Text = "0 B/s"
			return nil // skip first read
		}
		bps := float64(tx-sent) / now.Sub(lastTime).Seconds()
		lastTime = now
		sent = tx
		block.Text = fmtHuman(uint64(bps)) + "/s"
		return nil
	}, nil
}

func init() {
	Install("netspeed_rx", netspeedRx)
	Install("netspeed_tx", netspeedTx)
}
