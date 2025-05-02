package component

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

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
	return buf.String(), nil
}
