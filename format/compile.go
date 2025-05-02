package format

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/friedelschoen/st8/component"
	"github.com/friedelschoen/st8/notify"
)

type ComponentCall struct {
	Func component.Component
	Arg  string

	Length  int
	Padding string
	LeftPad bool
}

func literal(text string, _ *notify.Notification) (string, error) {
	return text, nil
}

var componentPattern = regexp.MustCompile(`^(\w+)(?:!(-)?([^1-9])?([0-9]+))?(?::(.*))?$`)

func parseComponentCall(text string, offset int) (ComponentCall, error) {
	m := componentPattern.FindStringSubmatch(text)
	if m == nil {
		return ComponentCall{}, fmt.Errorf("invalid component, beginning at %d: `%s`", offset, text)
	}

	name := m[1]
	padLeft := m[2] != ""
	padRune := m[3]
	padLength, _ := strconv.Atoi(m[4])
	arg := m[5]

	if padRune == "" {
		padRune = " "
	}

	compFunc, ok := component.Functions[name]
	if !ok {
		return ComponentCall{}, fmt.Errorf("undefined function, beginning at %d: %s", offset, name)
	}
	return ComponentCall{Func: compFunc, Arg: arg, LeftPad: padLeft, Length: padLength, Padding: padRune}, nil
}

func CompileFormat(input string) (ComponentFormat, error) {
	var calls ComponentFormat
	offset := 0

	input = strings.ReplaceAll(input, "\n", "")

	for {
		nextIdx := strings.Index(input, "${")
		if nextIdx == -1 {
			break
		}
		if nextIdx > 0 {
			calls = append(calls, ComponentCall{Func: literal, Arg: input[:nextIdx]})
			input = input[nextIdx:]
			offset += nextIdx
		}

		endIdx := strings.Index(input, "}")
		if endIdx == -1 {
			return nil, fmt.Errorf("unterminated call, beginning at %d: `%s`", offset, input)
		}

		call, err := parseComponentCall(input[2:endIdx], offset)
		if err != nil {
			return nil, err
		}
		calls = append(calls, call)
		input = input[endIdx+1:]
		offset += endIdx + 1
	}

	if len(input) > 0 {
		calls = append(calls, ComponentCall{Func: literal, Arg: input})
	}

	return calls, nil
}
