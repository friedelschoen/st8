package component

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/friedelschoen/st8/notify"
)

type BlockColor string

type Width string

type Alignment string

const (
	AlignLeft   Alignment = "left"
	AlignCenter Alignment = "center"
	AlignRight  Alignment = "right"
)

type Markup string

const (
	MarkupNone     Markup = "none"  // none
	MarkupPango    Markup = "pango" // pango
	MarkupMarkdown Markup = "none"  // none -> is converted to pango
)

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

type Block struct {
	// default text to display
	Text string `json:"full_text"`
	// text to display if not enough space
	Short string `json:"short_text,omitempty"`
	// color of this block
	TextColor BlockColor `json:"color,omitempty" conf:"text-color"`
	// color of this block
	BackgroundColor BlockColor `json:"background,omitempty" conf:"background-color"`
	// color of the border
	BorderColor BlockColor `json:"border,omitempty" conf:"border-color"`
	// Top border height in pixels (default 1).
	BorderTop int `json:"border_top,omitempty" conf:"border-top"`
	// Bottom border height in pixels (default 1).
	BorderBottom int `json:"border_bottom,omitempty" conf:"border-bottom"`
	// Left border width in pixels (default 1).
	BorderLeft int `json:"border_left,omitempty" conf:"border-left"`
	// Right border width in pixels (default 1).
	BorderRight int `json:"border_right,omitempty" conf:"border-right"`
	// minimum block width (either 10px or 10wh or a string which length will taken)
	Width Width `json:"min_width,omitempty" conf:"width"`
	// alignment
	Align Alignment `json:"align,omitempty" conf:"align"`
	// is urgent (result blinking)
	Urgent bool `json:"urgent,omitempty"`
	// put saperator after this block
	Seperate bool `json:"separator,omitempty" conf:"separate"`
	// separator width
	SeperatorWidth int `json:"separator_block_width,omitempty" conf:"separator-width"`
	// use markup
	Markup Markup `json:"markup,omitempty" conf:"markup"`
	// Identifier for click events
	Name string `json:"name,omitempty"`
	// Idenfifier
	ID int `json:"id"`
	// Event Handlers
	Handlers EventHandlers `json:"-"`
}

type Component func(block *Block, not *notify.Notification) error

type ComponentBuilder func(args map[string]string, handler *EventHandlers) (Component, error)

func (width Width) MarshalJSON() ([]byte, error) {
	switch {
	case strings.HasPrefix(string(width), "px"):
		cnt, err := strconv.Atoi(string(width)[:len(width)-2])
		if err != nil {
			return nil, err
		}
		return json.Marshal(cnt)
	case strings.HasPrefix(string(width), "wh"):
		cnt, err := strconv.Atoi(string(width)[:len(width)-2])
		if err != nil {
			return nil, err
		}
		return json.Marshal(strings.Repeat(" ", cnt))
	default:
		return json.Marshal(string(width))
	}
}

func Install(name string, builder ComponentBuilder) {
	Functions[name] = builder
}

var Functions = make(map[string]ComponentBuilder)

// fmtHuman formats bytes to a human-readable string, e.g. "1.4 GiB"
func fmtHuman(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := unit, 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// matches `text` against `pattern`, pattern may contain an asterisk to consume any string
func globMatch(pattern, text string) bool {
	prefix, suffix, ok := strings.Cut(pattern, "*")
	if ok {
		return strings.HasPrefix(text, prefix) && strings.HasSuffix(text, suffix)
	}
	return pattern == text
}
