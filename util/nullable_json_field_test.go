package util

import (
	"encoding/json"
	"reflect"
	"testing"
)

type TestStruct struct {
	Field      NullableJsonField[string] `json:"hello"`
	NullField  NullableJsonField[string] `json:"world"`
	ErrTrigger string                    `json:"err_trigger"`
}

func TestUnmarshalJSON_SimpleStruct(t *testing.T) {
	jsoned := `{"hello": "abc", "world": null, "err_trigger": "a"}`

	var s TestStruct
	err := json.Unmarshal([]byte(jsoned), &s)

	if err != nil {
		t.Fatalf("Unmarshalling valid json should not return an error")
	}

	expected := TestStruct{
		Field:      NullableJsonField[string]{Value: "abc", IsNull: false},
		NullField:  NullableJsonField[string]{Value: "", IsNull: true},
		ErrTrigger: "a",
	}

	if !reflect.DeepEqual(s, expected) {
		t.Fatalf("Unmarshalling json should return expected result")
	}
}

func TestUnmarshalJSON_InvalidSimpleStruct(t *testing.T) {
	jsoned := `{}`

	var s TestStruct
	err := json.Unmarshal([]byte(jsoned), &s)

	if err != nil {
		t.Fatalf("Unmarshalling invalid json SHOULD return an error")
	}
}

func TestMarshalJSON_SimpleStruct(t *testing.T) {
	marshalled := TestStruct{
		Field:      NullableJsonField[string]{Value: "abc", IsNull: false},
		NullField:  NullableJsonField[string]{Value: "", IsNull: true},
		ErrTrigger: "a",
	}

	jsoned, err := json.Marshal(marshalled)

	if err != nil {
		t.Fatalf("Marshalling valid json should not return an error")
	}

	expected := []byte(`{"hello":"abc","world":null,"err_trigger":"a"}`)

	if !reflect.DeepEqual(jsoned, expected) {
		t.Fatalf("Unmarshalling json should return expected result")
	}
}
