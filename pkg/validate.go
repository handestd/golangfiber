package pkg

import (
	"reflect"
)

func IsMapContainNil(s map[string]interface{}, omitKey []string) bool {
	for k, v := range s {
		if contains(omitKey, k) == false {
			if v == "" {
				return true
			}
		}
	}
	return false
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func IsStructContainNil(s interface{}) bool {
	val := reflect.ValueOf(s).Elem()

	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).IsNil() {
			return true
		}
	}
	return false
}
func CompareType(data1, data2 interface{}) bool {
	if reflect.TypeOf(data1) == reflect.TypeOf(data2) {
		return true
	}
	return false
}
