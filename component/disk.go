package component

import (
	"fmt"

	"github.com/friedelschoen/st8/notify"
	"golang.org/x/sys/unix"
)

func DiskFree(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		var stat unix.Statfs_t
		if err := unix.Statfs(args["path"], &stat); err != nil {
			return err
		}

		block.Text = fmtHuman(stat.Bavail * uint64(stat.Bsize))
		return nil
	}, nil
}

func DiskUsed(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		var stat unix.Statfs_t
		if err := unix.Statfs(args["path"], &stat); err != nil {
			return err
		}

		block.Text = fmtHuman((stat.Blocks - stat.Bfree) * uint64(stat.Bsize))
		return nil
	}, nil
}

func DiskTotal(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		var stat unix.Statfs_t
		if err := unix.Statfs(args["path"], &stat); err != nil {
			return err
		}

		block.Text = fmtHuman(stat.Bfree * uint64(stat.Bsize))
		return nil
	}, nil
}

func DiskPercentage(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		var stat unix.Statfs_t
		if err := unix.Statfs(args["path"], &stat); err != nil {
			return err
		}

		block.Text = fmt.Sprintf("%.0f", 100-(float64(stat.Bavail)/float64(stat.Blocks))*100)
		return nil
	}, nil
}
