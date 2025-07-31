package popupmenu

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type MenuItem struct {
	Id       string
	Text     string
	Children []MenuItem
}

func writeMenu(in io.Writer, items []MenuItem, depth int) {
	for _, item := range items {
		if item.Text == "" && item.Id == "" {
			fmt.Fprintln(in)
			continue
		}
		if item.Text == "" {
			item.Text = item.Id
		}
		fmt.Fprintf(in, "%s%s\t%s\n", strings.Repeat("\t", depth), item.Text, item.Id)
		writeMenu(in, item.Children, depth+1)
	}
}

func PopupMenu(items []MenuItem) (string, error) {
	cmd := exec.Command("xmenu")
	in, err := cmd.StdinPipe()
	var out strings.Builder
	cmd.Stdout = &out
	if err != nil {
		return "", err
	}
	if err := cmd.Start(); err != nil {
		return "", err
	}
	writeMenu(in, items, 0)
	in.Close()
	if err := cmd.Wait(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}
