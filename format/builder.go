package format

import (
	"errors"
	"strings"
	"sync"

	"github.com/friedelschoen/st8/notify"
)

type ComponentFormat []*ComponentCall

var ErrorString = "<error>"

func (cf ComponentFormat) Build(not *notify.Notification) (string, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error
	var length int
	results := make([]string, len(cf))

	for i, call := range cf {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := call.Func(call.Arg, not, &call.Cache)
			if err != nil {
				result = ErrorString
			}
			if call.Length != 0 && len(result) < call.Length {
				pad := strings.Repeat(call.Padding, call.Length-len(result))
				if call.LeftPad {
					result = result + pad
				} else {
					result = pad + result
				}
			}

			mu.Lock()
			defer mu.Unlock()
			results[i] = result
			length += len(result)
			errs = append(errs, err)
		}()
	}

	wg.Wait()

	bytes := make([]byte, 0, length)
	for _, s := range results {
		bytes = append(bytes, s...)
	}
	return string(bytes), errors.Join(errs...)
}
