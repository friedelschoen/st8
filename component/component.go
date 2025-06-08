package component

import (
	"fmt"
	"strings"

	"github.com/friedelschoen/st8/notify"
	"github.com/friedelschoen/st8/proto"
)

type Component func(block *proto.Block, not *notify.Notification) error

type ComponentBuilder func(args map[string]string, handler *proto.EventHandlers) (Component, error)

var Functions = make(map[string]ComponentBuilder)

func Install(name string, builder ComponentBuilder) {
	Functions[name] = builder
}

// fmtHuman formats bytes to a human-readable string, e.g. "1.4 GiB"
func fmtHuman(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := unit, 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// matches `text` against `pattern`, pattern may contain an asterisk to consume any string
func globMatch(pattern, text string) bool {
	prefix, suffix, ok := strings.Cut(pattern, "*")
	if ok {
		return strings.HasPrefix(text, prefix) && strings.HasSuffix(text, suffix)
	}
	return pattern == text
}
