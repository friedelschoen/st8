package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/friedelschoen/st8/format"
	"github.com/friedelschoen/st8/notify"
	"github.com/spf13/pflag"
)

var (
	statusFile     = pflag.StringP("status", "s", "", "path to status format")
	notifyFile     = pflag.StringP("notification", "n", "", "path to notification format")
	timeout        = pflag.Duration("not-timeout", 5*time.Second, "default timeout of a notification")
	rotateInterval = pflag.DurationP("rotate", "r", time.Second, "rotate notifications every ...")
	updateInterval = pflag.DurationP("update", "u", time.Second, "update interval")
	printFlag      = pflag.BoolP("print", "p", false, "print to stdout instead of using XStoreName")
	onceFlag       = pflag.BoolP("once", "1", false, "only print once (implies --print)")
	noWarn         = pflag.BoolP("no-warn", "w", false, "suppress command errors")
	helpFlag       = pflag.BoolP("help", "h", false, "show help and exit")
)

func readFormat(path string) format.ComponentFormat {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to read format file at %s: %v\n", path, err)
		os.Exit(1)
	}
	comp, err := format.CompileFormat(string(data))
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to parse format file at %s: %v\n", path, err)
		os.Exit(1)
	}
	return comp
}

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

	cStatus := readFormat(*statusFile)
	cNotify := readFormat(*notifyFile)

	runOnce := *onceFlag
	printMode := *printFlag || runOnce

	if runOnce {
		text, err := cStatus.Build(nil)
		if err != nil && !*noWarn {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Println(text)
		return
	}

	dpy := OpenDisplay()
	if dpy == nil {
		fmt.Fprintln(os.Stderr, "unable to open display")
		os.Exit(1)
	}
	defer dpy.Close()

	notifyChan, err := notify.NotifyStart()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to start daemon: %v", err)
		os.Exit(1)
	}
	defer notifyChan.Close()

	ticker := time.NewTicker(*updateInterval)
	defer ticker.Stop()

	var notifications []string
	var notifIndex int
	var notifTimer *time.Timer

	showNotif := func() {
		if len(notifications) == 0 {
			return
		}
		prefix := fmt.Sprintf("(%d/%d) ", notifIndex+1, len(notifications))
		text := prefix + notifications[notifIndex]
		if printMode {
			fmt.Println(text)
		} else {
			dpy.StoreName(text)
		}
	}

	rotateTicker := time.NewTicker(*rotateInterval)
	defer rotateTicker.Stop()

	for {
		select {
		case not := <-notifyChan.C:
			text, err := cNotify.Build(&not)
			if err != nil && !*noWarn {
				fmt.Fprintln(os.Stderr, err)
			}
			notifications = append(notifications, text)
			notifIndex = 0
			showNotif()

			if notifTimer != nil {
				notifTimer.Stop()
			}
			notifTimer = time.AfterFunc(*timeout, func() {
				notifications = nil
			})

		case <-rotateTicker.C:
			if len(notifications) > 1 {
				notifIndex = (notifIndex + 1) % len(notifications)
				showNotif()
			}

		case <-ticker.C:
			if len(notifications) == 0 {
				text, err := cStatus.Build(nil)
				if err != nil && !*noWarn {
					fmt.Fprintln(os.Stderr, err)
				}
				if printMode {
					fmt.Println(text)
				} else {
					dpy.StoreName(text)
				}
			}
		}
	}
}
