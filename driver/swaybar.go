//go:build swaybar
// +build swaybar

package driver

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/friedelschoen/st8/component"
	"github.com/friedelschoen/st8/driver/swaybarproto"
)

type SwayStatus struct {
	update chan<- struct{}
	enc    *json.Encoder
	dec    *json.Decoder

	handlers map[string]component.EventHandler
}

func init() {
	Drivers["swaybar"] = &SwayStatus{}
}

func (dpy *SwayStatus) Init(update chan<- struct{}) error {
	dpy.update = update
	dpy.handlers = make(map[string]component.EventHandler)
	dpy.enc = json.NewEncoder(os.Stdout)
	dpy.dec = json.NewDecoder(os.Stdin)
	hdr := swaybarproto.Header{
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
		var evt component.ClickEvent
		err := dpy.dec.Decode(&evt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to read event: %v", err)
			return
		}
		handler, ok := dpy.handlers[fmt.Sprintf("%s-%s", evt.Name, evt.Instance)]
		if !ok {
			continue
		}
		handler(evt)
		dpy.update <- struct{}{}
	}
}

func (dpy *SwayStatus) Close() error {
	return nil
}

func (dpy *SwayStatus) SetText(line []component.Block) error {
	body := make([]swaybarproto.SwayBlock, len(line))
	for i, block := range line {
		body[i].Block = block
		body[i].Instance = strconv.Itoa(i)
		if block.OnClick != nil {
			dpy.handlers[fmt.Sprintf("%s-%d", block.Name, i)] = block.OnClick
		}
	}
	err := dpy.enc.Encode(body)
	if err != nil {
		return err
	}
	fmt.Print(",")
	return nil
}
