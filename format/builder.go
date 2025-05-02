package format

import (
	"errors"
	"sync"

	"github.com/friedelschoen/st8/notify"
)

type ComponentFormat []ComponentCall

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
			result, err := call.Func(call.Arg, not)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errs = append(errs, err)
				result = ErrorString
			}
			results[i] = result
			length += len(result)
		}()
	}

	wg.Wait()

	bytes := make([]byte, 0, length)
	for _, s := range results {
		bytes = append(bytes, s...)
	}
	return string(bytes), errors.Join(errs...)
}
