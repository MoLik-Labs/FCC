package wirejson

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"unicode/utf16"
)

// EncodeASCII renders JSON like Python's json.dumps defaults: insertion-order
// Objects, ASCII escapes, and spaces after separators. Plain Go maps have no
// insertion order and are deliberately emitted in sorted-key order. Use
// Object recursively wherever byte-for-byte source ordering is required.
// Supported values are JSON primitives, Object, map[string]any and []any.
func EncodeASCII(value any) (string, error) {
	var output strings.Builder
	if err := writeASCII(&output, value); err != nil {
		return "", err
	}
	return output.String(), nil
}

func writeASCII(b *strings.Builder, value any) error {
	switch v := value.(type) {
	case nil:
		b.WriteString("null")
	case bool:
		b.WriteString(strconv.FormatBool(v))
	case string:
		writeASCIIString(b, v)
	case int:
		b.WriteString(strconv.Itoa(v))
	case int64:
		b.WriteString(strconv.FormatInt(v, 10))
	case float64:
		b.WriteString(pythonFloat(v))
	case json.Number:
		text := string(v)
		if len(text) == 0 || (text[0] != '-' && (text[0] < '0' || text[0] > '9')) || !json.Valid([]byte(text)) {
			return fmt.Errorf("invalid JSON number %q", text)
		}
		if strings.ContainsAny(text, ".eE") {
			f, err := strconv.ParseFloat(text, 64)
			if err != nil && !math.IsInf(f, 0) {
				return err
			}
			b.WriteString(pythonFloat(f))
		} else {
			// Python json.loads parses integer -0 as 0. All other valid
			// integer spellings retain arbitrary precision.
			if text == "-0" {
				text = "0"
			}
			b.WriteString(text)
		}
	case *Object:
		if v == nil {
			b.WriteString("null")
			break
		}
		b.WriteByte('{')
		for i, key := range v.Keys {
			if i > 0 {
				b.WriteString(", ")
			}
			writeASCIIString(b, key)
			b.WriteString(": ")
			if err := writeASCII(b, v.Fields[key]); err != nil {
				return err
			}
		}
		b.WriteByte('}')
	case Object:
		return writeASCII(b, &v)
	case map[string]any:
		keys := make([]string, 0, len(v))
		for key := range v {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		b.WriteByte('{')
		for i, key := range keys {
			if i > 0 {
				b.WriteString(", ")
			}
			writeASCIIString(b, key)
			b.WriteString(": ")
			if err := writeASCII(b, v[key]); err != nil {
				return err
			}
		}
		b.WriteByte('}')
	case []any:
		b.WriteByte('[')
		for i, item := range v {
			if i > 0 {
				b.WriteString(", ")
			}
			if err := writeASCII(b, item); err != nil {
				return err
			}
		}
		b.WriteByte(']')
	default:
		return fmt.Errorf("unsupported JSON type %T", value)
	}
	return nil
}

func pythonFloat(value float64) string {
	if math.IsNaN(value) {
		return "NaN"
	}
	if math.IsInf(value, 1) {
		return "Infinity"
	}
	if math.IsInf(value, -1) {
		return "-Infinity"
	}
	abs := math.Abs(value)
	if abs != 0 && (abs < 1e-4 || abs >= 1e16) {
		return strconv.FormatFloat(value, 'e', -1, 64)
	}
	text := strconv.FormatFloat(value, 'f', -1, 64)
	if !strings.Contains(text, ".") {
		text += ".0"
	}
	return text
}

// WriteASCIIString writes a Python json.dumps-compatible quoted JSON string.
func WriteASCIIString(b *strings.Builder, value string) {
	writeASCIIString(b, value)
}

func writeASCIIString(b *strings.Builder, value string) {
	b.WriteByte('"')
	for _, r := range value {
		switch r {
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		case '\b':
			b.WriteString(`\b`)
		case '\f':
			b.WriteString(`\f`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		default:
			if r < 0x20 || r >= 0x7f {
				if r <= 0xffff {
					writeHex16(b, uint16(r))
				} else {
					hi, lo := utf16.EncodeRune(r)
					writeHex16(b, uint16(hi))
					writeHex16(b, uint16(lo))
				}
			} else {
				b.WriteRune(r)
			}
		}
	}
	b.WriteByte('"')
}

func writeHex16(b *strings.Builder, v uint16) {
	b.WriteString(`\u`)
	b.WriteByte(hexChars[(v>>12)&15])
	b.WriteByte(hexChars[(v>>8)&15])
	b.WriteByte(hexChars[(v>>4)&15])
	b.WriteByte(hexChars[v&15])
}
