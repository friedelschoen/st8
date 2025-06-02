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
	enc *json.Encoder
	dec *json.Decoder
}

func init() {
	Drivers["swaybar"] = &SwayStatus{}
}

func (dpy *SwayStatus) Init() error {
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
	f, _ := os.Create("clickyevent.txt")
	defer f.Close()
	for {
		var evt swaybarproto.ClickEvent
		dpy.dec.Decode(&evt)

		bytes, _ := json.Marshal(evt)
		f.Write(bytes)
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
	}
	dpy.enc.Encode(body)
	fmt.Print(",")
	return nil
}
