package main

import (
	"fmt"
	"maps"
	"os"
	"path"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/friedelschoen/st8/driver"
	"github.com/friedelschoen/st8/format"
	"github.com/friedelschoen/st8/notify"
	"github.com/friedelschoen/st8/proto"
	"github.com/spf13/pflag"
)

func driverNames() string {
	keys := slices.Collect(maps.Keys(driver.Drivers))
	slices.Sort(keys)
	return strings.Join(keys, ", ")
}

var (
	statusFile     = pflag.StringP("status", "s", "", "path to status format")
	notifyFile     = pflag.StringP("notification", "n", "", "path to notification format")
	verify         = pflag.Bool("verify", false, "only verify config")
	timeout        = pflag.DurationP("notif-timeout", "N", 10*time.Second, "default timeout of a notification")
	rotateInterval = pflag.DurationP("rotate", "r", 2500*time.Millisecond, "rotate notifications every ...")
	updateInterval = pflag.DurationP("update", "u", time.Second, "update interval")
	driverFlag     = pflag.StringP("driver", "d", "stdout", "use driver: "+driverNames())
	onceFlag       = pflag.BoolP("once", "1", false, "only print once (implies --print)")
	quiet          = pflag.BoolP("quiet", "q", false, "suppress command errors")
	helpFlag       = pflag.BoolP("help", "h", false, "show help and exit")
	noNotify       = pflag.Bool("no-notify", false, "disable notifications")
)

func main() {
	pflag.Parse()

	if *helpFlag {
		pflag.Usage()
		os.Exit(0)
	}

	if *statusFile == "" || *notifyFile == "" {
		confdir, err := os.UserConfigDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to determite config-dir: %v\n", err)
			os.Exit(1)
		}
		if *statusFile == "" {
			*statusFile = path.Join(confdir, "st8", "status.txt")
		}
		if *notifyFile == "" {
			*notifyFile = path.Join(confdir, "st8", "notify.txt")
		}
	}

	cStatus, err := format.BuildComponents(*statusFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error in status-config: %v\n", err)
		os.Exit(1)
	}
	cNotify, err := format.BuildComponents(*notifyFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error in notify-config: %v\n", err)
		os.Exit(1)
	}

	if *verify {
		os.Exit(0)
	}

	drv, ok := driver.Drivers[*driverFlag]
	if !ok {
		fmt.Fprintf(os.Stdout, "not a valid driver: %s\n  valid drivers are: %s\n", *driverFlag, driverNames())
	}
	updateNow := make(chan struct{})
	if err := drv.Init(updateNow); err != nil {
		fmt.Fprintf(os.Stderr, "unable to initialize driver: %v\n", err)
		os.Exit(1)
	}
	defer drv.Close()

	if *onceFlag {
		text, err := cStatus.Build(nil)
		if err != nil && !*quiet {
			fmt.Fprintln(os.Stderr, err)
		}
		err = drv.SetText(text)
		if err != nil && !*quiet {
			fmt.Fprintln(os.Stderr, err)
		}
		return
	}

	notifyChannel := make(chan notify.Notification)
	var notifier *notify.NotificationDaemon

	if !*noNotify {
		notifier, err = notify.NotifyStart(notifyChannel)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to start daemon: %v", err)
			os.Exit(1)
		}
	}
	defer func() {
		if notifier != nil {
			notifier.Close()
		}
	}()

	var notifMu sync.Mutex
	var notifSet [][]proto.Block
	var notifIndex int

	showNotif := func() {
		if len(notifSet) == 0 {
			return
		}
		var prefix string
		if len(notifSet) != 1 {
			prefix = fmt.Sprintf("(%d/%d) ", notifIndex+1, len(notifSet))
		}
		blocks := []proto.Block{
			{Text: prefix},
		}
		blocks = append(blocks, notifSet[notifIndex]...)
		if err := drv.SetText(blocks); err != nil && !*quiet {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	updateTicker := time.NewTicker(*updateInterval)
	defer updateTicker.Stop()

	rotateTicker := time.NewTicker(*rotateInterval)
	defer rotateTicker.Stop()

	for {
		select {
		case not := <-notifyChannel:
			text, err := cNotify.Build(&not)
			if err != nil && !*quiet {
				fmt.Fprintln(os.Stderr, err)
			}
			notifMu.Lock()
			notifSet = append(notifSet, text)
			notifIndex = 0
			showNotif()
			notifMu.Unlock()

			nTimeout := *timeout
			if not.Timeout != 0 {
				nTimeout = not.Timeout
			}

			linesum := func(blocks []proto.Block) int {
				i := 0
				for _, blk := range blocks {
					i += blk.ID
				}
				return i
			}

			time.AfterFunc(nTimeout, func() {
				notifMu.Lock()
				defer notifMu.Unlock()
				notifSet = slices.DeleteFunc(notifSet, func(n []proto.Block) bool {
					return linesum(n) == linesum(text)
				})
			})

		case <-rotateTicker.C:
			notifMu.Lock()
			if len(notifSet) > 1 {
				showNotif()
			}
			notifMu.Unlock()

		case <-updateNow:
			if len(notifSet) == 0 {
				text, err := cStatus.Build(nil)
				if err != nil && !*quiet {
					fmt.Fprintln(os.Stderr, err)
				}
				if err := drv.SetText(text); err != nil && !*quiet {
					fmt.Fprintln(os.Stderr, err)
				}
			}

		case <-updateTicker.C:
			if len(notifSet) == 0 {
				text, err := cStatus.Build(nil)
				if err != nil && !*quiet {
					fmt.Fprintln(os.Stderr, err)
				}
				if err := drv.SetText(text); err != nil && !*quiet {
					fmt.Fprintln(os.Stderr, err)
				}
			}
		}
	}
}
