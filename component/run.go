package component

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/friedelschoen/st8/notify"
)

func RunCommand(cmdline string, _ *notify.Notification) (string, error) {
	var buf strings.Builder
	cmd := exec.Command("sh", "-c", cmdline)
	cmd.Stdin = nil
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("unable to execute `%s`: %w", cmdline, err)
	}
	return strings.TrimSpace(buf.String()), nil
}

type commandstate struct {
	output  string
	err     error
	time    time.Time
	running bool
}

var commandmutex sync.Mutex
var commandoutputs = make(map[string]commandstate)

func PeriodCommand(arg string, _ *notify.Notification) (string, error) {
	durstr, cmdline, ok := strings.Cut(arg, ",")
	if !ok {
		return "", fmt.Errorf("argument requires a comma")
	}
	dur, err := time.ParseDuration(durstr)
	if err != nil {
		return "", fmt.Errorf("invalid duration `%s`: %w", durstr, err)
	}

	commandmutex.Lock()
	defer commandmutex.Unlock()
	cache, ok := commandoutputs[cmdline]

	if !ok || (!cache.running && time.Since(cache.time) > dur) {
		cache.running = true
		commandoutputs[cmdline] = cache
		go func() {
			var buf strings.Builder
			cmd := exec.Command("sh", "-c", cmdline)
			cmd.Stdin = nil
			cmd.Stdout = &buf
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				cache.err = fmt.Errorf("unable to execute `%s`: %w", cmdline, err)
			}
			commandmutex.Lock()
			defer commandmutex.Unlock()
			cache.output = strings.TrimSpace(buf.String())
			cache.time = time.Now()
			cache.running = false
			commandoutputs[cmdline] = cache
		}()
	}

	return cache.output, cache.err
}
