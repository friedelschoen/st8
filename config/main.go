package config

import "time"

type MainConfig struct {
	Output         string        `conf:"driver.output"`
	Notifiers      string        `conf:"driver.notifiers"`
	StatusInterval time.Duration `conf:"status.interval"`
	NotifyTimeout  time.Duration `conf:"notification.timeout"`
	NotifyRotate   time.Duration `conf:"notification.rotate"`
	SocketNetwork  string        `conf:"socket-notifier.network"`
	SocketAddress  string        `conf:"socket-notifier.address"`
}

var DefaultConf = MainConfig{
	Output:         "stdout",
	Notifiers:      "", /* none */
	StatusInterval: 1 * time.Second,
	NotifyTimeout:  5 * time.Second,
	NotifyRotate:   1500 * time.Millisecond,
	SocketNetwork:  "tcp4",
	SocketAddress:  "0.0.0.0:4040",
}
