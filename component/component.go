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

type EventHandler func(evt ClickEvent)

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
	// put seperator after this block
	Seperate bool `json:"separator,omitempty" conf:"seperate"`
	// seperator width
	SeperatorWidth int `json:"separator_block_width,omitempty" conf:"seperator-width"`
	// use markup
	Markup Markup `json:"markup,omitempty" conf:"markup"`
	// Identifier for click events
	Name string `json:"name,omitempty"`
	// Click handler
	OnClick EventHandler `json:"-"`
	// Idenfifier
	ID int `json:"id"`
}

type Component func(block *Block, args map[string]string, not *notify.Notification, cache *any) error

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

var Functions = map[string]Component{
	"counter":           Counter,
	"battery_state":     BatteryState,
	"battery_perc":      BatteryPercentage,
	"battery_remaining": BatteryRemaining,
	"cat":               ReadFile,
	"cpu_perc":          CPUPercentage,
	"datetime":          Datetime,
	"disk_free":         DiskFree,
	"disk_perc":         DiskPercentage,
	"disk_total":        DiskTotal,
	"disk_used":         DiskUsed,
	"entropy":           EntropyAvailable,
	"gid":               Gid,
	"hostname":          Hostname,
	"ipv4":              IPv4,
	"ipv6":              IPv6,
	"kernel_release":    KernelRelease,
	"load_avg":          LoadAverage,
	"netspeed_rx":       NetspeedRx,
	"netspeed_tx":       NetspeedTx,
	"notify_appname":    NotifyAppName,
	"notify_appicon":    NotifyAppIcon,
	"notify_summary":    NotifySummary,
	"notify_body":       NotifyBody,
	"notify_actions":    NotifyActions,
	"period_command":    PeriodCommand,
	"num_files":         NumFiles,
	"ram_free":          RamFree,
	"ram_perc":          RamPercentage,
	"ram_total":         RamTotal,
	"ram_used":          RamUsed,
	"run_command":       RunCommand,
	"swap_free":         SwapFree,
	"swap_perc":         SwapPercentage,
	"swap_total":        SwapTotal,
	"swap_used":         SwapUsed,
	"temp":              Temperature,
	"uid":               Uid,
	"up":                Up,
	"uptime":            Uptime,
	"username":          Username,
	"wifi_essid":        WifiESSID,
	"wifi_perc":         WifiPerc,
}

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
