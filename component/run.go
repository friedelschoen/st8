package component

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/friedelschoen/st8/notify"
	"github.com/friedelschoen/st8/proto"
)

func runCommand(args map[string]string, events *proto.EventHandlers) (Component, error) {
	command, ok := args["command"]
	if !ok {
		return nil, fmt.Errorf("missing argument: command")
	}
	return func(block *proto.Block, not *notify.Notification) error {
		var buf strings.Builder
		cmd := exec.Command("sh", "-c", command)
		cmd.Stdin = nil
		cmd.Stdout = &buf
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("unable to execute `%s`: %w", command, err)
		}
		block.Text = strings.TrimSpace(buf.String())
		return nil
	}, nil
}

func periodCommand(args map[string]string, events *proto.EventHandlers) (Component, error) {
	var (
		output     = "?"
		commanderr error
		lastTime   time.Time
		running    bool
		mu         sync.Mutex
	)
	command, ok := args["command"]
	if !ok {
		return nil, fmt.Errorf("missing argument: command")
	}
	interval, ok := args["interval"]
	if !ok {
		return nil, fmt.Errorf("missing argument: interval")
	}

	dur, err := time.ParseDuration(interval)
	if err != nil {
		return nil, fmt.Errorf("invalid duration `%s`: %w", args["interval"], err)
	}

	return func(block *proto.Block, not *notify.Notification) error {
		mu.Lock()
		defer mu.Unlock()

		if !running && time.Since(lastTime) > dur {
			running = true
			go func() {
				var buf strings.Builder
				cmd := exec.Command("sh", "-c", command)
				cmd.Stdin = nil
				cmd.Stdout = &buf
				cmd.Stderr = os.Stderr

				if err := cmd.Run(); err != nil {
					commanderr = fmt.Errorf("unable to execute `%s`: %w", command, err)
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

func init() {
	Install("period_command", periodCommand)
	Install("run_command", runCommand)
}
