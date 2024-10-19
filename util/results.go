package util

type Result[T any] struct {
	Value T
	Error error
}

type IndexedResult[T any] struct {
	Value T
	Error error
	Index int
}
