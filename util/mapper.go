package util

import (
	"context"
	"sync"
	"time"
)

type Mappable[T any] interface {
	Items() []T
	Count() int
}

type MapperFunc[S, D any] func(ctx context.Context, s S) (D, error)

func AsyncListMapper[S any, A Mappable[S], D any](ctx context.Context, source A, mapFunc MapperFunc[S, D], timeout time.Duration) ([]D, []error) {
	sourceCount := len(source.Items())
	mapped := make([]D, sourceCount)
	errorsCollected := make([]error, 0)

	var wg sync.WaitGroup
	respCh := make(chan *IndexedResult[D], sourceCount)

	for i, v := range source.Items() {
		wg.Add(1)
		go func(index int, s S) {
			defer wg.Done()

			timeoutCtx, cancelTimeout := context.WithTimeout(ctx, timeout)
			defer cancelTimeout()

			if d, err := mapFunc(timeoutCtx, s); err != nil {
				respCh <- &IndexedResult[D]{
					Error: err,
					Index: index,
				}
			} else {
				respCh <- &IndexedResult[D]{
					Value: d,
					Index: i,
				}
			}
		}(i, v)
	}

	go func() {
		wg.Wait()
		close(respCh)
	}()

	for result := range respCh {
		if result.Error != nil {
			errorsCollected = append(errorsCollected, result.Error)
		} else {
			mapped[result.Index] = result.Value
		}
	}

	return mapped, errorsCollected
}
