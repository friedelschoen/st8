package format

import (
	"errors"
	"strings"
	"sync"

	"github.com/friedelschoen/st8/component"
	"github.com/friedelschoen/st8/notify"
)

type ComponentFormat []*ComponentCall

const ErrorString = "<error>"

var incrementer = 0

func (cf ComponentFormat) Build(not *notify.Notification) ([]component.Block, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error
	results := make([]component.Block, len(cf))

	for i, call := range cf {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := call.Block
			result.ID = incrementer
			incrementer++
			err := call.Func(&result, call.Arg, not, &call.Cache)
			if err != nil {
				result = call.Block
				result.Text = ErrorString
			}
			if call.Length != 0 && len(result.Text) < call.Length {
				pad := strings.Repeat(call.Padding, call.Length-len(result.Text))
				if call.LeftPad {
					result.Text += pad
				} else {
					result.Text = pad + result.Text
				}
			}
			result.Text = call.Prefix + result.Text + call.Suffix

			mu.Lock()
			defer mu.Unlock()
			results[i] = result
			errs = append(errs, err)
		}()
	}

	wg.Wait()

	return results, errors.Join(errs...)
}
