package proto

type ClickEvent struct {
	// The name of the block, if set.
	Name string `json:"name,omitempty"`
	// The instance of the block, if set.
	Instance string `json:"instance,omitempty"`
	// The absolute X location of the click.
	X int `json:"x"`
	// The absolute Y location of the click.
	Y int `json:"y"`
	// The X11 button number (or 0 if unmapped).
	Button int `json:"button"`
	// The event code corresponding to the button.
	Event int `json:"event"`
	// X position relative to the block’s top-left corner.
	RelativeX int `json:"relative_x"`
	// Y position relative to the block’s top-left corner.
	RelativeY int `json:"relative_y"`
	// Width of the block in pixels.
	Width int `json:"width"`
	// Height of the block in pixels.
	Height int `json:"height"`
}

type EventHandlers struct {
	OnClick func(evt ClickEvent)
}
