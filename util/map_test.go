package util

import (
	"reflect"
	"testing"
)

func TestSortedKeys(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"w": 54,
		"u": 12,
		"b": 189,
	}

	expected := []string{"a", "b", "u", "w"}

	if !reflect.DeepEqual(SortedKeys(m), expected) {
		t.Fatalf("Should return an alphabetically sorted array of the map keys")
	}
}
