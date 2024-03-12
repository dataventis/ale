package types

import (
	"bytes"
	"cmp"
	"fmt"
	"strings"

	"github.com/kode4food/comb/basics"
)

type (
	// Record describes an Object that allows a fixed set of fields,
	// each of which has a keyword
	Record struct {
		basic
		fields
	}

	// Field describes one of the fields of a RecordType
	Field struct {
		Value Type
		Name  string
	}

	fields []Field
)

// MakeRecord declares a new RecordType that only allows a fixed set of Field
// entries, each being identified by a BasicKeyword and having a specified Type
func MakeRecord(fields ...Field) *Record {
	return &Record{
		basic:  BasicObject,
		fields: fields,
	}
}

func (r *Record) Fields() []Field {
	return r.fields
}

func (r *Record) Name() string {
	return fmt.Sprintf("record(%s)", r.name())
}

func (r *Record) Accepts(c *Checker, other Type) bool {
	if r == other {
		return true
	}
	if other, ok := other.(*Record); ok {
		rf := r.fields
		of := other.Fields()
		if len(rf) > len(of) {
			return false
		}
		om := fields(of).toMap()
		for k, v := range rf.toMap() {
			if tv, ok := om[k]; !ok || !c.AcceptsChild(v, tv) {
				return false
			}
		}
		return true
	}
	return false
}

func (r *Record) Equal(other Type) bool {
	if r == other {
		return true
	}
	if other, ok := other.(*Record); ok {
		return r.basic.Equal(other.basic) && r.equal(other.fields)
	}
	return false
}

func (f fields) toMap() map[string]Type {
	res := map[string]Type{}
	for _, p := range f {
		res[p.Name] = p.Value
	}
	return res
}

func (f fields) sorted() fields {
	return basics.SortFunc(f, func(l, r Field) int {
		return cmp.Compare(l.Name, r.Name)
	})
}

func (f fields) name() string {
	var buf bytes.Buffer
	for i, elem := range f.sorted() {
		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteRune('"')
		buf.WriteString(strings.ReplaceAll(elem.Name, `"`, `\"`))
		buf.WriteString("\"->")
		buf.WriteString(elem.Value.Name())
	}
	return buf.String()
}

func (f fields) equal(other fields) bool {
	if len(f) != len(other) {
		return false
	}
	fs := f.sorted()
	os := other.sorted()
	for i, l := range fs {
		r := os[i]
		if l.Name != r.Name || !l.Value.Equal(r.Value) {
			return false
		}
	}
	return true
}
