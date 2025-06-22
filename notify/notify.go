package notify

import (
	"io"
	"time"

	"github.com/friedelschoen/st8/config"
)

type Notification struct {
	AppName string        `json:"name"`
	AppIcon string        `json:"icon"`
	Summary string        `json:"summary"`
	Body    string        `json:"body"`
	Actions []string      `json:"actions"`
	Timeout time.Duration `json:"timeout"`
}

type NotificationDaemon func(*config.MainConfig, chan<- Notification) (io.Closer, error)

var Functions = make(map[string]NotificationDaemon)

func Install(name string, daemon NotificationDaemon) {
	Functions[name] = daemon
}
