// Code generated by "stringer -type DataType"; DO NOT EDIT.

package desc

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[UNKNOWN-0]
	_ = x[STRING-1]
	_ = x[NUMBER-2]
	_ = x[BOOLEAN-3]
}

const _DataType_name = "UNKNOWNSTRINGNUMBERBOOLEAN"

var _DataType_index = [...]uint8{0, 7, 13, 19, 26}

func (i DataType) String() string {
	if i < 0 || i >= DataType(len(_DataType_index)-1) {
		return "DataType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _DataType_name[_DataType_index[i]:_DataType_index[i+1]]
}
