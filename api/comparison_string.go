// Code generated by "stringer -type=Comparison -linecomment"; DO NOT EDIT.

package api

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[LessThan - -1]
	_ = x[EqualTo-0]
	_ = x[GreaterThan-1]
	_ = x[Incomparable-2]
}

const _Comparison_name = "Less ThanEqual ToGreater ThanNot Comparable"

var _Comparison_index = [...]uint8{0, 9, 17, 29, 43}

func (i Comparison) String() string {
	i -= -1
	if i < 0 || i >= Comparison(len(_Comparison_index)-1) {
		return "Comparison(" + strconv.FormatInt(int64(i+-1), 10) + ")"
	}
	return _Comparison_name[_Comparison_index[i]:_Comparison_index[i+1]]
}
