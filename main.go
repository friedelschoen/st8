package main

import (
	"fmt"
	"io"
	"maps"
	"os"
	"path"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/friedelschoen/st8/component"
	"github.com/friedelschoen/st8/config"
	"github.com/friedelschoen/st8/driver"
	"github.com/friedelschoen/st8/format"
	"github.com/friedelschoen/st8/notify"
	"github.com/friedelschoen/st8/proto"
	"github.com/spf13/pflag"
)

func driverNames[T any](drvs map[string]T) string {
	keys := slices.Collect(maps.Keys(drvs))
	slices.Sort(keys)
	return strings.Join(keys, ", ")
}

func defaultConfigPath() string {
	confdir, err := os.UserConfigDir()
	if err != nil {
		homedir, err := os.UserHomeDir()
		if err != nil {
			homedir = "/root"
		}
		confdir = path.Join(homedir, ".config")
	}
	return path.Join(confdir, "st8")
}

var (
	configPath    = pflag.StringP("config", "c", defaultConfigPath(), "path to config-directory")
	verify        = pflag.Bool("verify", false, "only verify config")
	driverFlag    = pflag.StringP("output", "T", "", "output to, available drivers: "+driverNames(driver.Drivers))
	notifiersFlag = pflag.StringP("notifier", "n", "", "enable notifiers (delimited by comma), available drivers: "+driverNames(notify.Functions))
	onceFlag      = pflag.BoolP("once", "1", false, "only print once")
	quiet         = pflag.BoolP("quiet", "q", false, "suppress command errors")
	helpFlag      = pflag.BoolP("help", "h", false, "show help and exit")
)

func parseConfig(conf *config.MainConfig, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return nil /* ignoring */
	}
	defer file.Close()

	for section, values := range config.ParseConfig(file, filename) {
		err := config.UnmarshalConf(values, section, conf)
		if err != nil {
			return err
		}
	}

	return nil
}

func procHooks() error {
	hookdir, _ := os.ReadDir(path.Join(*configPath, "hooks"))
	for _, entry := range hookdir {
		component.Install(entry.Name(), component.HookComponent(path.Join(*configPath, "hooks", entry.Name())))
	}
	return nil
}

func main() {
	pflag.Parse()

	if *helpFlag {
		pflag.Usage()
		os.Exit(0)
	}

	conf := config.DefaultConf
	err := parseConfig(&conf, path.Join(*configPath, "config.ini"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error in config: %v\n", err)
		os.Exit(1)
	}

	err = procHooks()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cStatus, err := format.BuildComponents(path.Join(*configPath, "status.ini"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error in status-format: %v\n", err)
		os.Exit(1)
	}
	cNotify, err := format.BuildComponents(path.Join(*configPath, "notification.ini"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error in notification-format: %v\n", err)
		os.Exit(1)
	}

	if *verify {
		os.Exit(0)
	}

	if *driverFlag != "" {
		conf.Output = *driverFlag
	}

	if *notifiersFlag != "" {
		conf.Notifiers = *notifiersFlag
	}

	drv, ok := driver.Drivers[conf.Output]
	if !ok {
		fmt.Fprintf(os.Stdout, "not a valid driver: %s\n  valid drivers are: %s\n", conf.Output, driverNames(driver.Drivers))
		os.Exit(0)
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
	var notifyClosers []io.Closer
	for drvname := range strings.SplitSeq(conf.Notifiers, ",") {
		drvname = strings.TrimSpace(drvname)
		if drvname == "" {
			continue
		}
		daemon, ok := notify.Functions[drvname]
		if !ok {
			fmt.Fprintf(os.Stdout, "not a valid driver: %s\n  valid drivers are: %s\n", drvname, driverNames(notify.Functions))
			os.Exit(1)
		}
		closer, err := daemon(&conf, notifyChannel)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to initialize driver: %v\n", err)
			os.Exit(1)
		}
		notifyClosers = append(notifyClosers, closer)
	}

	defer func() {
		for _, notifier := range notifyClosers {
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

	updateTicker := time.NewTicker(conf.StatusInterval)
	defer updateTicker.Stop()

	rotateTicker := time.NewTicker(conf.NotifyRotate)
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

			nTimeout := conf.NotifyTimeout
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
