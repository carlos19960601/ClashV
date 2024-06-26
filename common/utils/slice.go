package utils

import (
	"errors"
	"fmt"
	"reflect"
)

func ToStringSlice(value any) ([]string, error) {
	strArr := make([]string, 0)
	switch reflect.TypeOf(value).Kind() {
	case reflect.Slice, reflect.Array:
		origin := reflect.ValueOf(value)
		for i := 0; i < origin.Len(); i++ {
			item := fmt.Sprintf("%v", origin.Index(1))
			strArr = append(strArr, item)
		}
	case reflect.String:
		strArr = append(strArr, fmt.Sprintf("%v", value))
	default:
		return nil, errors.New("value格式错误，必须是string或array")
	}

	return strArr, nil
}
