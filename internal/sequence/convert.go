package sequence

import (
	"bytes"

	"github.com/kode4food/ale/data"
)

// ToList takes any sequence and converts it to a List
func ToList(s data.Sequence) data.Sequence {
	switch s := s.(type) {
	case *data.List:
		return s
	case data.CountedSequence:
		res := make(data.Values, s.Count())
		idx := 0
		for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
			res[idx] = f
			idx++
		}
		return data.NewList(res...)
	default:
		return uncountedToList(s)
	}
}

func uncountedToList(s data.Sequence) data.Sequence {
	return data.NewList(uncountedToValues(s)...)
}

// ToValues takes any sequence and converts it to a value array
func ToValues(s data.Sequence) data.Values {
	switch s := s.(type) {
	case data.ValuerSequence:
		return s.Values()
	case data.CountedSequence:
		res := make(data.Values, s.Count())
		idx := 0
		for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
			res[idx] = f
			idx++
		}
		return res
	default:
		return uncountedToValues(s)
	}
}

func uncountedToValues(s data.Sequence) data.Values {
	res := data.Values{}
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		res = append(res, f)
	}
	return res
}

// ToVector takes any sequence and converts it to a vector
func ToVector(s data.Sequence) data.Vector {
	switch s := s.(type) {
	case data.Vector:
		return s
	default:
		v := ToValues(s)
		return data.NewVector(v...)
	}
}

// ToObject takes any sequence and converts it to an Object
func ToObject(s data.Sequence) (*data.Object, error) {
	switch s := s.(type) {
	case *data.Object:
		return s, nil
	default:
		v := ToValues(s)
		return data.ValuesToObject(v...)
	}
}

// ToStr takes any sequence and attempts to convert it to a String
func ToStr(s data.Sequence) data.String {
	if st, ok := s.(data.String); ok {
		return st
	}
	var buf bytes.Buffer
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		if f == data.Null {
			continue
		}
		buf.WriteString(f.String())
	}
	return data.String(buf.String())
}
