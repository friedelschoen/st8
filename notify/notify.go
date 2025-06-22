package notify

import (
	"io"
	"time"
)

type Notification struct {
	AppName string
	AppIcon string
	Summary string
	Body    string
	Actions []string
	Timeout time.Duration
}

type NotificationDaemon func(chan<- Notification) (io.Closer, error)

var Functions = make(map[string]NotificationDaemon)

func Install(name string, daemon NotificationDaemon) {
	Functions[name] = daemon
}
