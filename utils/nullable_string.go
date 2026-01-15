package utils

import (
	"encoding/json"
)

type NullableString struct {
	IsSet  bool
	IsNull bool
	Value  string
}

func NewNullableString(value string) *NullableString {
	return &NullableString{
		IsSet:  true,
		IsNull: false,
		Value:  value,
	}
}

func NullableStringNull() *NullableString {
	return &NullableString{
		IsSet:  true,
		IsNull: true,
	}
}

func (n *NullableString) MarshalJSON() ([]byte, error) {
	if !n.IsSet {
		return []byte("undefined"), nil
	}

	if n.IsNull {
		return []byte("null"), nil
	}

	return json.Marshal(&n.Value)
}

func (n *NullableString) UnmarshalJSON(bytes []byte) error {
	n.IsSet = true

	if string(bytes) == "null" {
		n.IsNull = true
		return nil
	}

	var s string
	if err := json.Unmarshal(bytes, &s); err != nil {
		return err
	}
	n.Value = s
	return nil
}

var _ json.Marshaler = &NullableString{}
var _ json.Unmarshaler = &NullableString{}
