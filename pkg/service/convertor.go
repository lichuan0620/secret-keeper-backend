package service

import (
	"reflect"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// Converter is used to convert string to specific type.
type Converter func([]string) (interface{}, error)

var converters = map[reflect.Type]Converter{
	reflect.TypeOf(bool(false)):        ConvertToBool,
	reflect.TypeOf(int(0)):             ConvertToInt,
	reflect.TypeOf(int8(0)):            ConvertToInt8,
	reflect.TypeOf(int16(0)):           ConvertToInt16,
	reflect.TypeOf(int32(0)):           ConvertToInt32,
	reflect.TypeOf(int64(0)):           ConvertToInt64,
	reflect.TypeOf(uint(0)):            ConvertToUint,
	reflect.TypeOf(uint8(0)):           ConvertToUint8,
	reflect.TypeOf(uint16(0)):          ConvertToUint16,
	reflect.TypeOf(uint32(0)):          ConvertToUint32,
	reflect.TypeOf(uint64(0)):          ConvertToUint64,
	reflect.TypeOf(float32(0)):         ConvertToFloat32,
	reflect.TypeOf(float64(0)):         ConvertToFloat64,
	reflect.TypeOf(string("")):         ConvertToString,
	reflect.TypeOf(time.Time{}):        ConvertToTime,
	reflect.TypeOf(time.Duration(0)):   ConvertToDuration,
	reflect.TypeOf(new(bool)):          ConvertToBoolP,
	reflect.TypeOf(new(int)):           ConvertToIntP,
	reflect.TypeOf(new(int8)):          ConvertToInt8P,
	reflect.TypeOf(new(int16)):         ConvertToInt16P,
	reflect.TypeOf(new(int32)):         ConvertToInt32P,
	reflect.TypeOf(new(int64)):         ConvertToInt64P,
	reflect.TypeOf(new(uint)):          ConvertToUintP,
	reflect.TypeOf(new(uint8)):         ConvertToUint8P,
	reflect.TypeOf(new(uint16)):        ConvertToUint16P,
	reflect.TypeOf(new(uint32)):        ConvertToUint32P,
	reflect.TypeOf(new(uint64)):        ConvertToUint64P,
	reflect.TypeOf(new(float32)):       ConvertToFloat32P,
	reflect.TypeOf(new(float64)):       ConvertToFloat64P,
	reflect.TypeOf(new(string)):        ConvertToStringP,
	reflect.TypeOf(new(time.Time)):     ConvertToTimeP,
	reflect.TypeOf(new(time.Duration)): ConvertToDurationP,
	reflect.TypeOf([]bool{}):           ConvertToBoolSlice,
	reflect.TypeOf([]int{}):            ConvertToIntSlice,
	reflect.TypeOf([]float64{}):        ConvertToFloat64Slice,
	reflect.TypeOf([]string{}):         ConvertToStringSlice,
}

// ConverterFor gets Converter for specified type.
func ConverterFor(typ reflect.Type) (Converter, error) {
	ret := converters[typ]
	if ret == nil {
		return nil, errors.Errorf("cannot convert string to %v", typ)
	}
	return ret, nil
}

// MustConverterFor is similar to ConverterFor but does not check for error.
func MustConverterFor(typ reflect.Type) Converter {
	return converters[typ]
}

func buildConvertError(v, t string) error {
	return errors.Errorf("cannot convert \"%s\" to type %s", v, t)
}

// ConvertToBool converts []string to bool.
func ConvertToBool(data []string) (interface{}, error) {
	target, err := strconv.ParseBool(data[0])
	if err != nil {
		return nil, buildConvertError(data[0], "bool")
	}
	return target, nil
}

// ConvertToBoolP converts []string to *bool.
func ConvertToBoolP(data []string) (interface{}, error) {
	target, err := strconv.ParseBool(data[0])
	if err != nil {
		return nil, buildConvertError(data[0], "bool")
	}
	return &target, nil
}

// ConvertToInt converts []string to int.
func ConvertToInt(data []string) (interface{}, error) {
	target, err := strconv.ParseInt(data[0], 10, 0)
	if err != nil {
		return nil, buildConvertError(data[0], "int")
	}
	return int(target), nil
}

// ConvertToIntP converts []string to *int.
func ConvertToIntP(data []string) (interface{}, error) {
	target, err := strconv.ParseInt(data[0], 10, 0)
	if err != nil {
		return nil, buildConvertError(data[0], "int")
	}
	value := int(target)
	return &value, nil
}

// ConvertToInt8 converts []string to int8.
func ConvertToInt8(data []string) (interface{}, error) {
	target, err := strconv.ParseInt(data[0], 10, 8)
	if err != nil {
		return nil, buildConvertError(data[0], "int8")
	}
	return int8(target), nil
}

// ConvertToInt8P converts []string to *int8.
func ConvertToInt8P(data []string) (interface{}, error) {
	target, err := strconv.ParseInt(data[0], 10, 8)
	if err != nil {
		return nil, buildConvertError(data[0], "int8")
	}
	value := int8(target)
	return &value, nil
}

// ConvertToInt16 converts []string to int16.
func ConvertToInt16(data []string) (interface{}, error) {
	target, err := strconv.ParseInt(data[0], 10, 16)
	if err != nil {
		return nil, buildConvertError(data[0], "int16")
	}
	return int16(target), nil
}

// ConvertToInt16P converts []string to *int16.
func ConvertToInt16P(data []string) (interface{}, error) {
	target, err := strconv.ParseInt(data[0], 10, 16)
	if err != nil {
		return nil, buildConvertError(data[0], "int16")
	}
	value := int16(target)
	return &value, nil
}

// ConvertToInt32 converts []string to int32.
func ConvertToInt32(data []string) (interface{}, error) {
	target, err := strconv.ParseInt(data[0], 10, 32)
	if err != nil {
		return nil, buildConvertError(data[0], "int32")
	}
	return int32(target), nil
}

// ConvertToInt32P converts []string to *int32.
func ConvertToInt32P(data []string) (interface{}, error) {
	target, err := strconv.ParseInt(data[0], 10, 32)
	if err != nil {
		return nil, buildConvertError(data[0], "int32")
	}
	value := int32(target)
	return &value, nil
}

// ConvertToInt64 converts []string to int64.
func ConvertToInt64(data []string) (interface{}, error) {
	target, err := strconv.ParseInt(data[0], 10, 64)
	if err != nil {
		return nil, buildConvertError(data[0], "int64")
	}
	return target, nil
}

// ConvertToInt64P converts []string to *int64.
func ConvertToInt64P(data []string) (interface{}, error) {
	target, err := strconv.ParseInt(data[0], 10, 64)
	if err != nil {
		return nil, buildConvertError(data[0], "int64")
	}
	return &target, nil
}

// ConvertToUint converts []string to uint.
func ConvertToUint(data []string) (interface{}, error) {
	target, err := strconv.ParseUint(data[0], 10, 0)
	if err != nil {
		return nil, buildConvertError(data[0], "uint")
	}
	return uint(target), nil
}

// ConvertToUintP converts []string to *uint.
func ConvertToUintP(data []string) (interface{}, error) {
	target, err := strconv.ParseUint(data[0], 10, 0)
	if err != nil {
		return nil, buildConvertError(data[0], "uint")
	}
	value := uint(target)
	return &value, nil
}

// ConvertToUint8 converts []string to uint8.
func ConvertToUint8(data []string) (interface{}, error) {
	target, err := strconv.ParseUint(data[0], 10, 8)
	if err != nil {
		return nil, buildConvertError(data[0], "uint8")
	}
	return uint8(target), nil
}

// ConvertToUint8P converts []string to *uint8.
func ConvertToUint8P(data []string) (interface{}, error) {
	target, err := strconv.ParseUint(data[0], 10, 8)
	if err != nil {
		return nil, buildConvertError(data[0], "uint8")
	}
	value := uint8(target)
	return &value, nil
}

// ConvertToUint16 converts []string to uint16.
func ConvertToUint16(data []string) (interface{}, error) {
	target, err := strconv.ParseUint(data[0], 10, 16)
	if err != nil {
		return nil, buildConvertError(data[0], "uint16")
	}
	return uint16(target), nil
}

// ConvertToUint16P converts []string to *uint16.
func ConvertToUint16P(data []string) (interface{}, error) {
	target, err := strconv.ParseUint(data[0], 10, 16)
	if err != nil {
		return nil, buildConvertError(data[0], "uint16")
	}
	value := uint16(target)
	return &value, nil
}

// ConvertToUint32 converts []string to uint32.
func ConvertToUint32(data []string) (interface{}, error) {
	target, err := strconv.ParseUint(data[0], 10, 32)
	if err != nil {
		return nil, buildConvertError(data[0], "uint32")
	}
	return uint32(target), nil
}

// ConvertToUint32P converts []string to *uint32.
func ConvertToUint32P(data []string) (interface{}, error) {
	target, err := strconv.ParseUint(data[0], 10, 32)
	if err != nil {
		return nil, buildConvertError(data[0], "uint32")
	}
	value := uint32(target)
	return &value, nil
}

// ConvertToUint64 converts []string to uint64.
func ConvertToUint64(data []string) (interface{}, error) {
	target, err := strconv.ParseUint(data[0], 10, 64)
	if err != nil {
		return nil, buildConvertError(data[0], "uint64")
	}
	return target, nil
}

// ConvertToUint64P converts []string to *uint64.
func ConvertToUint64P(data []string) (interface{}, error) {
	target, err := strconv.ParseUint(data[0], 10, 64)
	if err != nil {
		return nil, buildConvertError(data[0], "uint64")
	}
	return &target, nil
}

// ConvertToFloat32 converts []string to float32.
func ConvertToFloat32(data []string) (interface{}, error) {
	target, err := strconv.ParseFloat(data[0], 32)
	if err != nil {
		return nil, buildConvertError(data[0], "float32")
	}
	return float32(target), nil
}

// ConvertToFloat32P converts []string to *float32.
func ConvertToFloat32P(data []string) (interface{}, error) {
	target, err := strconv.ParseFloat(data[0], 32)
	if err != nil {
		return nil, buildConvertError(data[0], "float32")
	}
	value := float32(target)
	return &value, nil
}

// ConvertToFloat64 converts []string to float64.
func ConvertToFloat64(data []string) (interface{}, error) {
	target, err := strconv.ParseFloat(data[0], 64)
	if err != nil {
		return nil, buildConvertError(data[0], "float64")
	}
	return target, nil
}

// ConvertToFloat64P converts []string to *float64.
func ConvertToFloat64P(data []string) (interface{}, error) {
	target, err := strconv.ParseFloat(data[0], 64)
	if err != nil {
		return nil, buildConvertError(data[0], "float64")
	}
	return &target, nil
}

// ConvertToTime converts []string to time.Time.
func ConvertToTime(data []string) (interface{}, error) {
	target, err := time.Parse(time.RFC3339, data[0])
	if err != nil {
		return nil, buildConvertError(data[0], "Time")
	}
	return target, nil
}

// ConvertToTimeP return the first element's pointer in []string.
func ConvertToTimeP(data []string) (interface{}, error) {
	target, err := time.Parse(time.RFC3339, data[0])
	if err != nil {
		return nil, buildConvertError(data[0], "Time")
	}
	return &target, nil
}

// ConvertToDuration converts []string to time.Duration.
func ConvertToDuration(data []string) (interface{}, error) {
	target, err := time.ParseDuration(data[0])
	if err != nil {
		return nil, buildConvertError(data[0], "Duration")
	}
	return target, nil
}

// ConvertToDurationP return the first element's pointer in []string.
func ConvertToDurationP(data []string) (interface{}, error) {
	target, err := time.ParseDuration(data[0])
	if err != nil {
		return nil, buildConvertError(data[0], "Duration")
	}
	return &target, nil
}

// ConvertToString return the first element in []string.
func ConvertToString(data []string) (interface{}, error) {
	return data[0], nil
}

// ConvertToStringP return the first element's pointer in []string.
func ConvertToStringP(data []string) (interface{}, error) {
	return &data[0], nil
}

// ConvertToBoolSlice converts all elements in data to bool, and return []bool
func ConvertToBoolSlice(data []string) (interface{}, error) {
	ret := make([]bool, len(data))
	for i := range data {
		r, err := ConvertToBool(data[i : i+1])
		if err != nil {
			return nil, err
		}
		ret[i] = r.(bool)
	}
	return ret, nil
}

// ConvertToIntSlice converts all elements in data to int, and return []int
func ConvertToIntSlice(data []string) (interface{}, error) {
	ret := make([]int, len(data))
	for i := range data {
		r, err := ConvertToInt(data[i : i+1])
		if err != nil {
			return nil, err
		}
		ret[i] = r.(int)
	}
	return ret, nil
}

// ConvertToFloat64Slice converts all elements in data to float64, and return []float64
func ConvertToFloat64Slice(data []string) (interface{}, error) {
	ret := make([]float64, len(data))
	for i := range data {
		r, err := ConvertToFloat64(data[i : i+1])
		if err != nil {
			return nil, err
		}
		ret[i] = r.(float64)
	}
	return ret, nil
}

// ConvertToStringSlice return all strings in data.
func ConvertToStringSlice(data []string) (interface{}, error) {
	return data, nil
}
