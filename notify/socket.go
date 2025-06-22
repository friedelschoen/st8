package notify

import (
	"encoding/json"
	"io"
	"net"

	"github.com/friedelschoen/st8/config"
)

func handleConn(stream io.ReadWriteCloser, channel chan<- Notification) {
	dec := json.NewDecoder(stream)
	var not Notification
	for {
		err := dec.Decode(&not)
		if err != nil {
			return
		}
		channel <- not
	}
}

func acceptSocket(listener net.Listener, channel chan<- Notification) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		go handleConn(conn, channel)
	}
}

func startSocketDaemon(conf *config.MainConfig, channel chan<- Notification) (io.Closer, error) {
	listener, err := net.Listen(conf.SocketNetwork, conf.SocketAddress)
	if err != nil {
		return nil, err
	}
	go acceptSocket(listener, channel)
	return listener, nil
}

func init() {
	Install("socket", startSocketDaemon)
}
