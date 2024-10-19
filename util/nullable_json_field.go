package util

import (
	"bytes"
	"encoding/json"
)

type NullableJsonField[T any] struct {
	Value  T
	IsNull bool
}

func (njs *NullableJsonField[T]) UnmarshalJSON(data []byte) error {
	str := string(data)

	if str == `null` {
		njs.IsNull = true
		return nil
	}

	njs.IsNull = false

	return json.Unmarshal(data, &njs.Value)
}

func (njs NullableJsonField[T]) MarshalJSON() ([]byte, error) {
	var data bytes.Buffer

	if njs.IsNull {
		data.WriteString(`null`)
		return data.Bytes(), nil
	}

	return json.Marshal(njs.Value)
}
