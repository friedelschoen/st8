package proto

import (
	"encoding/json"
	"strconv"
	"strings"
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
