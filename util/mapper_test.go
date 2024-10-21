package util

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"
)

type Bookself[T string] struct {
	books []T
}

func (b Bookself[T]) Items() []T {
	return b.books
}

func (b Bookself[T]) Count() int {
	return len(b.books)
}

func TestAsyncListMapper(t *testing.T) {
	source := Bookself[string]{
		books: []string{"lotr", "the hobbit", "d&d manual"},
	}
	ctx := context.Background()

	expected := []string{"LOTR", "THE HOBBIT", "D&D MANUAL"}

	actual, errorsCollected := AsyncListMapper(ctx, source, func(ctx context.Context, el string) (string, error) {
		return strings.ToUpper(el), nil
	}, time.Hour)

	if len(errorsCollected) != 0 {
		t.Fatalf("Async list mapping should not have collected errors")
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Async list mapping should have returned expected result")
	}
}
