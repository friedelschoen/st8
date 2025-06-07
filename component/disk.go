package component

import (
	"fmt"

	"github.com/friedelschoen/st8/notify"
	"golang.org/x/sys/unix"
)

func diskFree(args map[string]string, events *EventHandlers) (Component, error) {
	path, ok := args["path"]
	if !ok {
		return nil, fmt.Errorf("missing argument: path")
	}
	return func(block *Block, not *notify.Notification) error {
		var stat unix.Statfs_t
		if err := unix.Statfs(path, &stat); err != nil {
			return err
		}

		block.Text = fmtHuman(stat.Bavail * uint64(stat.Bsize))
		return nil
	}, nil
}

func diskUsed(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		var stat unix.Statfs_t
		if err := unix.Statfs(args["path"], &stat); err != nil {
			return err
		}

		block.Text = fmtHuman((stat.Blocks - stat.Bfree) * uint64(stat.Bsize))
		return nil
	}, nil
}

func diskTotal(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		var stat unix.Statfs_t
		if err := unix.Statfs(args["path"], &stat); err != nil {
			return err
		}

		block.Text = fmtHuman(stat.Bfree * uint64(stat.Bsize))
		return nil
	}, nil
}

func diskPercentage(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		var stat unix.Statfs_t
		if err := unix.Statfs(args["path"], &stat); err != nil {
			return err
		}

		block.Text = fmt.Sprintf("%.0f", 100-(float64(stat.Bavail)/float64(stat.Blocks))*100)
		return nil
	}, nil
}

func init() {
	Install("disk_free", diskFree)
	Install("disk_perc", diskPercentage)
	Install("disk_total", diskTotal)
	Install("disk_used", diskUsed)
}
