package format

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/friedelschoen/st8/component"
	"github.com/friedelschoen/st8/config"
	"github.com/friedelschoen/st8/proto"
)

type Format struct {
	Length         int
	Padding        rune
	Align          Alignment
	Prefix, Suffix string
}

func (format *Format) Do(text string) string {
	if len(text) >= format.Length {
		return format.Prefix + text + format.Suffix
	}
	var leftPad, rightPad int
	space := format.Length - len(text)
	switch format.Align {
	case AlignLeft:
		rightPad = space
	case AlignCenter:
		leftPad = space / 2
		rightPad = space / 2
		if space%2 != 0 {
			rightPad++
		}
	case AlignRight:
		leftPad = space
	}
	var out strings.Builder
	out.WriteString(format.Prefix)
	for range leftPad {
		out.WriteRune(format.Padding)
	}
	out.WriteString(text)
	for range rightPad {
		out.WriteRune(format.Padding)
	}
	out.WriteString(format.Suffix)

	return out.String()
}

type ComponentCall struct {
	Func         component.Component
	Handlers     proto.EventHandlers
	DefaultBlock proto.Block

	Format      Format
	ShortFormat Format
}

type Alignment int

const (
	AlignLeft Alignment = iota
	AlignCenter
	AlignRight
)

var componentPattern = regexp.MustCompile(`^(?:([<^>-])?([^1-9])?([0-9]+))?$`)

func parseFormat(text string) (format Format, err error) {
	begin := strings.IndexByte(text, '{')
	if begin == -1 {
		return format, fmt.Errorf("format does not contain {}: %s", text)
	}
	format.Prefix = text[:begin]
	text = text[begin+1:]

	end := strings.IndexByte(text, '}')
	if end == -1 {
		return format, fmt.Errorf("unmatched `}`: %s", text)
	}
	format.Suffix = text[end+1:]
	text = text[:end]

	m := componentPattern.FindStringSubmatch(text)
	if m == nil {
		return format, fmt.Errorf("invalid format: %s", text)
	}
	chr, _ := utf8.DecodeRuneInString(m[1])
	switch chr {
	case '-', '<':
		format.Align = AlignLeft
	case '^':
		format.Align = AlignCenter
	case '>', utf8.RuneError:
		format.Align = AlignRight
	}
	chr, _ = utf8.DecodeRuneInString(m[2])
	if chr == utf8.RuneError {
		format.Padding = ' '
	} else {
		format.Padding = chr
	}
	format.Length, _ = strconv.Atoi(m[3])

	return
}

func BuildComponents(filename string) (ComponentFormat, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var result ComponentFormat
	for compname, values := range config.ParseConfig(file, filename) {
		call := &ComponentCall{}
		builder, ok := component.Functions[compname]
		if !ok {
			return nil, fmt.Errorf("unknown component: %s", compname)
		}
		var err error
		call.Func, err = builder(values, &call.Handlers)
		if err != nil {
			return nil, err
		}

		if format, ok := values["format"]; ok {
			call.Format, err = parseFormat(format)
			if err != nil {
				return nil, err
			}
		}
		if shortformat, ok := values["short_format"]; ok {
			call.ShortFormat, err = parseFormat(shortformat)
			if err != nil {
				return nil, err
			}
		}

		call.DefaultBlock.Name = compname
		if err := config.UnmarshalConf(values, "", &call.DefaultBlock); err != nil {
			return nil, err
		}
		result = append(result, call)
	}

	return result, nil
}
