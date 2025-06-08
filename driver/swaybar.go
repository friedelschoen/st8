//go:build swaybar
// +build swaybar

package driver

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/friedelschoen/st8/proto"
)

type Header struct {
	Version        int            `json:"version"`                // The protocol version to use. Currently, this must be 1
	ClickEvents    bool           `json:"click_events,omitempty"` // Whether to receive click event information to standard input
	ContinueSignal syscall.Signal `json:"cont_signal,omitempty"`  // The signal that swaybar should send to continue processing
	StopSignal     syscall.Signal `json:"stop_signal,omitempty"`  // The signal that swaybar should send to stop processing
}

type SwayBlock struct {
	proto.Block
	Instance string `json:"instance,omitempty"` // Secondary identifier for click events.
}

type SwayStatus struct {
	update chan<- struct{}
	enc    *json.Encoder
	dec    *json.Decoder

	handlers map[string]proto.EventHandlers
}

func init() {
	Drivers["swaybar"] = &SwayStatus{}
}

func (dpy *SwayStatus) Init(update chan<- struct{}) error {
	dpy.update = update
	dpy.handlers = make(map[string]proto.EventHandlers)
	dpy.enc = json.NewEncoder(os.Stdout)
	dpy.dec = json.NewDecoder(os.Stdin)
	hdr := Header{
		Version:     1,
		ClickEvents: true,
	}
	dpy.enc.Encode(hdr)
	fmt.Print("[")

	go dpy.eventLoop()

	return nil
}

func (dpy *SwayStatus) eventLoop() {
	if tok, err := dpy.dec.Token(); err != nil || tok != json.Delim('[') {
		return
	}
	for {
		var evt proto.ClickEvent
		err := dpy.dec.Decode(&evt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to read event: %v", err)
			return
		}
		handler, ok := dpy.handlers[fmt.Sprintf("%s-%s", evt.Name, evt.Instance)]
		if !ok {
			continue

		}
		if handler.OnClick == nil {
			continue
		}
		handler.OnClick(evt)
		dpy.update <- struct{}{}
	}
}

func (dpy *SwayStatus) Close() error {
	return nil
}

func (dpy *SwayStatus) SetText(line []proto.Block) error {
	body := make([]SwayBlock, len(line))
	for i, block := range line {
		body[i].Block = block
		body[i].Instance = strconv.Itoa(i)
		dpy.handlers[fmt.Sprintf("%s-%d", block.Name, i)] = block.Handlers
	}
	err := dpy.enc.Encode(body)
	if err != nil {
		return err
	}
	fmt.Print(",")
	return nil
}
