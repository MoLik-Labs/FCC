package wirejson

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// DecodeOrdered parses one JSON value, preserving object insertion order and
// number spelling. Like Python's JSON decoder, duplicate keys replace their
// values without changing the key's original position.
func DecodeOrdered(raw string) (any, error) {
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.UseNumber()
	value, err := readOrdered(decoder)
	if err != nil {
		return nil, err
	}
	if _, err = decoder.Token(); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("trailing JSON value")
		}
		return nil, err
	}
	return value, nil
}

func readOrdered(decoder *json.Decoder) (any, error) {
	token, err := decoder.Token()
	if err != nil {
		return nil, err
	}
	delim, isDelim := token.(json.Delim)
	if !isDelim {
		return token, nil
	}
	switch delim {
	case '{':
		obj := NewObject(0)
		for decoder.More() {
			key, err := decoder.Token()
			if err != nil {
				return nil, err
			}
			name, ok := key.(string)
			if !ok {
				return nil, fmt.Errorf("expected JSON object key")
			}
			value, err := readOrdered(decoder)
			if err != nil {
				return nil, err
			}
			obj.Set(name, value)
		}
		if _, err := decoder.Token(); err != nil {
			return nil, err
		}
		return obj, nil
	case '[':
		arr := []any{}
		for decoder.More() {
			value, err := readOrdered(decoder)
			if err != nil {
				return nil, err
			}
			arr = append(arr, value)
		}
		if _, err := decoder.Token(); err != nil {
			return nil, err
		}
		return arr, nil
	default:
		return nil, fmt.Errorf("unexpected JSON delimiter %q", delim)
	}
}
