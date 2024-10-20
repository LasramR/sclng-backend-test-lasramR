package util

// Utility type for multi value channels
type Result[T any] struct {
	Value T
	Error error
}

// Utility type for multi value buffered channels
type IndexedResult[T any] struct {
	Value T
	Error *error
	Index int
}
