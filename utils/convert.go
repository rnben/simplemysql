package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func ToInt64(i interface{}) (v int64) {
	switch t := i.(type) {
	case []uint8:
		v, _ = strconv.ParseInt(UI8ToA(i), 10, 64)
	case int:
		v = int64(t)
	case string:
		v, _ = strconv.ParseInt(t, 10, 64)
	case int64:
		v = t
	case int32:
		v = int64(t)
	case int8:
		v = int64(t)
	case int16:
		v = int64(t)
	case uint:
		v = int64(t)
	case uint8:
		v = int64(t)
	case uint16:
		v = int64(t)
	case uint32:
		v = int64(t)
	case uint64:
		v = int64(t)
	case float64:
		v = int64(t)

	}

	return v
}

func ToInt(i interface{}) (v int) {
	switch t := i.(type) {
	case int:
		v = t
	case int8:
		v = int(t)
	case int16:
		v = int(t)
	case int32:
		v = int(t)
	case int64:
		v = int(t)
	case uint:
		v = int(t)
	case uint8:
		v = int(t)
	case uint16:
		v = int(t)
	case uint32:
		v = int(t)
	case uint64:
		v = int(t)
	case float64:
		v = int(t)
	case string:
		vv, _ := strconv.ParseInt(t, 10, 64)
		v = int(vv)
	case []uint8:
		vv, _ := strconv.ParseInt(UI8ToA(i), 10, 64)
		v = int(vv)
	}

	return v
}

func ToString(i interface{}) (v string) {
	switch t := i.(type) {
	case string:
		v = t
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		v = fmt.Sprintf("%d", t)
	case float32, float64:
		v = fmt.Sprintf("%v", t)
	case []uint8:
		v = UI8ToA(t)
	}

	return v
}

func ToFloat64(i interface{}) (v float64) {
	switch t := i.(type) {
	case float32:
		v = float64(t)
	case float64:
		v = t
	case string:
		v, _ = strconv.ParseFloat(t, 64)
	case int, int8, int16, int32, int64:
		v = float64(ToInt64(i))
	case []byte:
		v, _ = strconv.ParseFloat(string(i.([]byte)), 64)
	}

	return v
}

func UI8ToA(i interface{}) string {
	if v, ok := i.(string); ok {
		return v
	}

	return string(UI8ToB(i))
}

func UI8ToB(i interface{}) (b []byte) {
	if v, ok := i.([]uint8); ok {
		b = append(b, v...)
	}

	return b
}

func StringMapConvert(fs map[string]interface{}, ts reflect.Type) interface{} {
	switch ts.Elem().Kind() {
	case reflect.String:
		res := make(map[string]string, len(fs))
		for k, v := range fs {
			res[k] = v.(string)
		}

		return res
	case reflect.Bool:
		res := make(map[string]bool, len(fs))
		for k, v := range fs {
			res[k] = v.(bool)
		}

		return res
	case reflect.Int:
		res := make(map[string]int, len(fs))
		for k, v := range fs {
			res[k] = ToInt(v)
		}

		return res
	case reflect.Int64:
		res := make(map[string]int64, len(fs))
		for k, v := range fs {
			res[k] = ToInt64(v)
		}

		return res
	case reflect.Float64:
		res := make(map[string]float64, len(fs))
		for k, v := range fs {
			res[k] = ToFloat64(v)
		}

		return res
	default:
		return fs
	}
}

func SliceInterfaceConvert(fs []interface{}, ts reflect.Type) interface{} {
	switch ts.Elem().Kind() {
	case reflect.String:
		res := make([]string, len(fs))
		for i, v := range fs {
			res[i] = v.(string)
		}

		return res
	case reflect.Bool:
		res := make([]bool, len(fs))
		for i, v := range fs {
			res[i] = v.(bool)
		}

		return res
	case reflect.Int:
		res := make([]int, len(fs))
		for i, v := range fs {
			res[i] = ToInt(v)
		}

		return res
	case reflect.Int64:
		res := make([]int64, len(fs))
		for i, v := range fs {
			res[i] = ToInt64(v)
		}

		return res
	case reflect.Float64:
		res := make([]float64, len(fs))
		for i, v := range fs {
			res[i] = ToFloat64(v)
		}

		return res
	default:
		return fs
	}
}

func ToBool(i interface{}) bool {
	switch t := i.(type) {
	case string:
		return strings.EqualFold(t, "true")
	case int, int8, int16, int32, int64:
		return ToInt64(i) > 0
	case bool:
		return t
	}

	return false
}

func Quote(s string) string {
	return "`" + s + "`"
}
