package format

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/friedelschoen/st8/component"
)

type ComponentCall struct {
	Func  component.Component
	Arg   map[string]string
	Block component.Block

	Length  int
	Padding string
	LeftPad bool

	Cache any

	Prefix, Suffix string
}

var componentPattern = regexp.MustCompile(`^(?:(-)?([^1-9])?([0-9]+))?$`)

func parseConfig(file io.Reader, filename string) iter.Seq2[string, map[string]string] {
	return func(yield func(string, map[string]string) bool) {
		scan := bufio.NewScanner(file)
		current := make(map[string]string)
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
				if len(section) > 0 && !yield(section, current) {
					return
				}

				section = newsection
				current = make(map[string]string)
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
		var ok bool
		call.Func, ok = component.Functions[compname]
		if !ok {
			return nil, fmt.Errorf("unknown component: %s", compname)
		}
		call.Arg = values

		if format, ok := values["format"]; ok {
			begin := strings.IndexByte(format, '{')
			if begin == -1 {
				return nil, fmt.Errorf("in component `%s`: format does not contain {}: %s", compname, format)
			}
			call.Prefix = format[:begin]
			format = format[begin+1:]

			end := strings.IndexByte(format, '}')
			if end == -1 {
				return nil, fmt.Errorf("in component `%s`: unmatched `}`: %s", compname, format)
			}
			call.Suffix = format[end+1:]
			format = format[:end]

			m := componentPattern.FindStringSubmatch(format)
			if m == nil {
				return nil, fmt.Errorf("in component `%s`: invalid format: %s", compname, format)
			}
			call.LeftPad = m[1] != ""
			call.Padding = m[2]
			call.Length, _ = strconv.Atoi(m[3])
		}
		if len(call.Padding) == 0 {
			call.Padding = " "
		}

		call.Block.Name = compname
		if err := UnmarshalConf(values, &call.Block); err != nil {
			return nil, err
		}
		result = append(result, call)
	}

	return result, nil
}
