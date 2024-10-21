package util

import (
	"context"
	"sync"
	"time"
)

// Ensure that type T can yield inner items and count
type Mappable[T any] interface {
	Items() []T
	Count() int
}

// Mapping function of an AsyncListMapper operation
type MapperFunc[S, D any] func(ctx context.Context, s S) (D, error)

// Given a Mappable source of type S, asynchronously transform the element to an array of type D
// Also returns a list of errors that occured during mapping process, check for error with len(errs) != 0
func AsyncListMapper[S any, A Mappable[S], D any](ctx context.Context, source A, mapFunc MapperFunc[S, D], timeout time.Duration) ([]D, []*error) {
	sourceCount := len(source.Items())
	mapped := make([]D, sourceCount)
	errorsCollected := make([]*error, 0)

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
					Error: &err,
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
