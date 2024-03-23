package utils

import (
	"strconv"
)

func StringToInt64(sValue string, defValue int64) int64 {
	value := defValue
	iValue, err := strconv.ParseInt(sValue, 0, 64)
	if err == nil {
		value = int64(iValue)
	}
	return value
}

func StringToInt32(sValue string, defValue int32) int32 {
	value := defValue
	iValue, err := strconv.ParseInt(sValue, 0, 32)
	if err == nil {
		value = int32(iValue)
	}
	return value
}

func StringToInt(sValue string, defValue int) int {
	value := defValue
	iValue, err := strconv.ParseInt(sValue, 0, 32)
	if err == nil {
		value = int(iValue)
	}
	return value
}

func StringToFloat32(sValue string, defValue float32) float32 {
	value := defValue
	iValue, err := strconv.ParseFloat(sValue, 0)
	if err == nil {
		value = float32(iValue)
	}
	return value
}

func StringToFloat64(sValue string, defValue float64) float64 {
	value := defValue
	iValue, err := strconv.ParseFloat(sValue, 64)
	if err == nil {
		value = iValue
	}
	return value
}
