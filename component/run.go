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

func RunCommand(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	var buf strings.Builder
	cmd := exec.Command("sh", "-c", args["command"])
	cmd.Stdin = nil
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unable to execute `%s`: %w", args["command"], err)
	}
	block.Text = strings.TrimSpace(buf.String())
	return nil
}

type commandstate struct {
	output  string
	err     error
	time    time.Time
	running bool
	mu      sync.Mutex
}

func PeriodCommand(block *Block, args map[string]string, not *notify.Notification, cacheptr *any) error {
	durstr, cmdline, ok := strings.Cut(args["command"], ",")
	if !ok {
		return fmt.Errorf("argument requires a comma")
	}
	dur, err := time.ParseDuration(durstr)
	if err != nil {
		return fmt.Errorf("invalid duration `%s`: %w", durstr, err)
	}

	/* commandstate pointer to avoid copying mutexes */
	var cache *commandstate
	if *cacheptr != nil {
		cache = (*cacheptr).(*commandstate)
	}
	if cache == nil {
		cache = new(commandstate)
		*cacheptr = cache
	}
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if !cache.running && time.Since(cache.time) > dur {
		cache.running = true
		go func() {
			var buf strings.Builder
			cmd := exec.Command("sh", "-c", cmdline)
			cmd.Stdin = nil
			cmd.Stdout = &buf
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				cache.err = fmt.Errorf("unable to execute `%s`: %w", cmdline, err)
			}
			cache.mu.Lock()
			defer cache.mu.Unlock()
			cache.output = strings.TrimSpace(buf.String())
			cache.time = time.Now()
			cache.running = false
		}()
	}

	block.Text = cache.output
	return cache.err
}
