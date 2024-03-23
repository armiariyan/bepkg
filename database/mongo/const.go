package mongo

import "errors"

//ErrInvalidTimeType Error if invalid time type happened
var ErrInvalidTimeType = errors.New("Invalid time type")

//ErrOutputRowsNotSlice Error output value is not slice
var ErrOutputRowsNotSlice = errors.New("Output rows is not a slice")

//ErrOutputValNotPointer output value is not a pointer message
var ErrOutputValNotPointer = errors.New("OutputVal is not a pointer")

//ErrOutputValNotSlicePointer slice pointer
var ErrOutputValNotSlicePointer = errors.New("OutputVal is not a pointer of slice")
