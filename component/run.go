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

func RunCommand(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
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
	}, nil
}

func PeriodCommand(args map[string]string, events *EventHandlers) (Component, error) {
	var (
		output     = "?"
		commanderr error
		lastTime   time.Time
		running    bool
		mu         sync.Mutex
	)

	dur, err := time.ParseDuration(args["interval"])
	if err != nil {
		return nil, fmt.Errorf("invalid duration `%s`: %w", args["interval"], err)
	}

	return func(block *Block, not *notify.Notification) error {
		mu.Lock()
		defer mu.Unlock()

		if !running && time.Since(lastTime) > dur {
			running = true
			go func() {
				var buf strings.Builder
				cmd := exec.Command("sh", "-c", args["command"])
				cmd.Stdin = nil
				cmd.Stdout = &buf
				cmd.Stderr = os.Stderr

				if err := cmd.Run(); err != nil {
					commanderr = fmt.Errorf("unable to execute `%s`: %w", args["command"], err)
				}
				mu.Lock()
				defer mu.Unlock()
				output = strings.TrimSpace(buf.String())
				lastTime = time.Now()
				running = false
			}()
		}

		block.Text = output
		return commanderr
	}, nil
}
