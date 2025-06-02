package swaybarproto

import "github.com/friedelschoen/st8/component"

type SwayBlock struct {
	component.Block
	Instance string `json:"instance,omitempty"` // Secondary identifier for click events.
}
