package format

import (
	"fmt"
	"strings"

	"github.com/friedelschoen/st8/component"
)

type ComponentCall struct {
	Func component.Component
	Arg  string
}

func literal(text string) (string, error) {
	return text, nil
}

func parseComponentCall(text string, offset int) (ComponentCall, error) {
	name, arg, _ := strings.Cut(text, ":")
	compFunc, ok := component.Functions[name]
	if !ok {
		return ComponentCall{}, fmt.Errorf("undefined function, beginning at %d: %s", offset, name)
	}
	return ComponentCall{Func: compFunc, Arg: arg}, nil
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
