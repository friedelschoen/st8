package format

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"maps"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/friedelschoen/st8/component"
	"github.com/friedelschoen/st8/proto"
)

type Format struct {
	Length         int
	Padding        rune
	Align          Alignment
	Prefix, Suffix string
}

func (format *Format) Do(text string) string {
	if format.Length == 0 || len(text) >= format.Length {
		return text
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

var componentPattern = regexp.MustCompile(`^(?:([<^>-])?([^1-9])?([0-9]+))?$`)

func parseConfig(file io.Reader, filename string) iter.Seq2[string, map[string]string] {
	return func(yield func(string, map[string]string) bool) {
		scan := bufio.NewScanner(file)
		current := make(map[string]string)
		var base map[string]string
		var section string
		var linenr int
		for scan.Scan() {
			line := scan.Text()
			linenr++

			if idx := strings.IndexAny(line, ";#"); idx != -1 {
				line = line[:idx]
			}
			line = strings.TrimSpace(line)
			if len(line) == 0 {
				continue
			}

			if line[0] == '[' {
				end := strings.IndexByte(line, ']')
				if end != len(line)-1 {
					fmt.Fprintf(os.Stderr, "%s:%d: garbage found after `]`: %s\n", filename, linenr, line[end:])
				}
				newsection := strings.TrimSpace(line[1:end])
				if len(newsection) == 0 {
					fmt.Fprintf(os.Stderr, "%s:%d: section is empty\n", filename, linenr)
					continue
				}
				if section == "" {
					base = current
				} else if !yield(section, current) {
					return
				}

				section = newsection
				current = maps.Clone(base)
				continue
			}

			key, value, ok := strings.Cut(line, "=")
			if !ok {
				fmt.Fprintf(os.Stderr, "%s:%d: not a key-value pair: %s\n", filename, linenr, line)
				continue
			}
			key = strings.TrimSpace(key)
			if len(key) == 0 {
				fmt.Fprintf(os.Stderr, "%s:%d: key is empty\n", filename, linenr)
				continue
			}
			value = strings.TrimSpace(value)
			if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
				value = value[1 : len(value)-1]
			}
			current[key] = value
		}
		if len(section) > 0 {
			yield(section, current)
		}
	}
}

func BuildComponents(filename string) (ComponentFormat, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var result ComponentFormat
	for compname, values := range parseConfig(file, filename) {
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
		if err := UnmarshalConf(values, &call.DefaultBlock); err != nil {
			return nil, err
		}
		result = append(result, call)
	}

	return result, nil
}
