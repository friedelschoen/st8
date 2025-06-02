package swaybarproto

type ClickEvent struct {
	Name      string `json:"name,omitempty"`     // The name of the block, if set.
	Instance  string `json:"instance,omitempty"` // The instance of the block, if set.
	X         int    `json:"x"`                  // The absolute X location of the click.
	Y         int    `json:"y"`                  // The absolute Y location of the click.
	Button    int    `json:"button"`             // The X11 button number (or 0 if unmapped).
	Event     int    `json:"event"`              // The event code corresponding to the button.
	RelativeX int    `json:"relative_x"`         // X position relative to the block’s top-left corner.
	RelativeY int    `json:"relative_y"`         // Y position relative to the block’s top-left corner.
	Width     int    `json:"width"`              // Width of the block in pixels.
	Height    int    `json:"height"`             // Height of the block in pixels.
}
