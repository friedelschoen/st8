package format

import (
	"errors"
	"sync"

	"github.com/friedelschoen/st8/notify"
	"github.com/friedelschoen/st8/proto"
)

type ComponentFormat []*ComponentCall

const ErrorString = "<error>"

var incrementer = 0

func (cf ComponentFormat) Build(not *notify.Notification) ([]proto.Block, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error
	results := make([]proto.Block, len(cf))

	for i, call := range cf {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := call.DefaultBlock
			result.Handlers = call.Handlers
			result.ID = incrementer
			incrementer++
			err := call.Func(&result, not)
			if err != nil {
				result = call.DefaultBlock
				result.Text = ErrorString
			}
			if result.Short == "" {
				result.Short = result.Text
			}
			result.Text = call.Format.Do(result.Text)
			result.Short = call.Format.Do(result.Short)

			mu.Lock()
			defer mu.Unlock()
			results[i] = result
			errs = append(errs, err)
		}()
	}

	wg.Wait()

	return results, errors.Join(errs...)
}
