package main

import (
	"fmt"
	"os"
	"time"

	"github.com/friedelschoen/st8/format"
	"github.com/spf13/pflag"
)

var (
	fileFlag     = pflag.StringP("file", "f", "", "use format file instead of argument")
	printFlag    = pflag.BoolP("print", "p", false, "print to stdout instead of using XStoreName")
	intervalFlag = pflag.DurationP("interval", "i", time.Second, "interval to update (default 1s)")
	notifyTime   = pflag.DurationP("notification-time", "n", 5*time.Second, "duration to show notification")
	onceFlag     = pflag.BoolP("once", "1", false, "only print once (implies --print)")
	noWarn       = pflag.BoolP("no-warn", "w", false, "suppress command errors")
	helpFlag     = pflag.BoolP("help", "h", false, "show help and exit")
)

func main() {
	pflag.Parse()

	if *helpFlag {
		pflag.Usage()
		os.Exit(0)
	}

	var formatStr string
	if *fileFlag != "" {
		data, err := os.ReadFile(*fileFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to read format file: %v\n", err)
			os.Exit(1)
		}
		formatStr = string(data)
	} else if pflag.NArg() > 0 {
		formatStr = pflag.Arg(0)
	} else {
		fmt.Fprintln(os.Stderr, "format string or --file required")
		os.Exit(1)
	}

	cf, err := format.CompileFormat(formatStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to parse format: %v\n", err)
		os.Exit(1)
	}

	runOnce := *onceFlag
	printMode := *printFlag || runOnce

	if runOnce {
		text, err := cf.Build()
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

	notify, err := NotifyStart()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to start daemon: %v", err)
		os.Exit(1)
	}
	defer notify.Close()

	ticker := time.NewTicker(*intervalFlag)
	defer ticker.Stop()

	for {
		select {
		case not := <-notify.C:
			if printMode {
				fmt.Println(not.summary)
			} else {
				dpy.StoreName(not.summary)
			}
			time.Sleep(*notifyTime)

		case <-ticker.C:
			text, err := cf.Build()
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
