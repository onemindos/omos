// Package jcs produces the JSON Canonicalization Scheme form (RFC 8785) of a JSON
// document. asset-command.md §5 and §6.3 define params_hash and the command signature
// over this form.
//
// Covered: object members sorted by UTF-16 code units, no insignificant whitespace,
// strings with only the escapes RFC 8785 §3.2.2.2 requires, and numbers in the
// ECMAScript Number.prototype.toString form (via encoding/json's float64 encoder,
// which follows ES6 formatting, with -0 written as 0). Numbers are IEEE 754 doubles,
// as RFC 8785 requires.
package jcs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"unicode/utf16"
	"unicode/utf8"
)

// Canonicalize returns the RFC 8785 form of the JSON document in data.
func Canonicalize(data []byte) ([]byte, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}
	if dec.More() {
		return nil, fmt.Errorf("jcs: trailing data")
	}
	var buf bytes.Buffer
	if err := write(&buf, v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Marshal JSON-encodes v and canonicalizes the result.
func Marshal(v any) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return Canonicalize(b)
}

func write(buf *bytes.Buffer, v any) error {
	switch t := v.(type) {
	case nil:
		buf.WriteString("null")
	case bool:
		buf.WriteString(strconv.FormatBool(t))
	case json.Number:
		f, err := strconv.ParseFloat(string(t), 64)
		if err != nil || math.IsInf(f, 0) || math.IsNaN(f) {
			return fmt.Errorf("jcs: number %q not representable", t)
		}
		if f == 0 {
			f = 0 // drop the sign of -0
		}
		b, err := json.Marshal(f)
		if err != nil {
			return err
		}
		buf.Write(b)
	case string:
		writeString(buf, t)
	case []any:
		buf.WriteByte('[')
		for i, e := range t {
			if i > 0 {
				buf.WriteByte(',')
			}
			if err := write(buf, e); err != nil {
				return err
			}
		}
		buf.WriteByte(']')
	case map[string]any:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool { return lessUTF16(keys[i], keys[j]) })
		buf.WriteByte('{')
		for i, k := range keys {
			if i > 0 {
				buf.WriteByte(',')
			}
			writeString(buf, k)
			buf.WriteByte(':')
			if err := write(buf, t[k]); err != nil {
				return err
			}
		}
		buf.WriteByte('}')
	default:
		return fmt.Errorf("jcs: unexpected type %T", v)
	}
	return nil
}

func lessUTF16(a, b string) bool {
	ua, ub := utf16.Encode([]rune(a)), utf16.Encode([]rune(b))
	for i := 0; i < len(ua) && i < len(ub); i++ {
		if ua[i] != ub[i] {
			return ua[i] < ub[i]
		}
	}
	return len(ua) < len(ub)
}

func writeString(buf *bytes.Buffer, s string) {
	buf.WriteByte('"')
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		i += size
		switch r {
		case '"':
			buf.WriteString(`\"`)
		case '\\':
			buf.WriteString(`\\`)
		case '\b':
			buf.WriteString(`\b`)
		case '\f':
			buf.WriteString(`\f`)
		case '\n':
			buf.WriteString(`\n`)
		case '\r':
			buf.WriteString(`\r`)
		case '\t':
			buf.WriteString(`\t`)
		default:
			if r < 0x20 {
				fmt.Fprintf(buf, `\u%04x`, r)
			} else {
				buf.WriteRune(r)
			}
		}
	}
	buf.WriteByte('"')
}
