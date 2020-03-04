package utils

import (
	"reflect"
	"strings"
)

type StructToMapCallback func(item interface{}, fields []string) map[string]interface{}

func StructToMap(item interface{}, fields []string) map[string]interface{} {
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	mapped := map[string]interface{}{}

	for _, field := range fields {
		for i := 0; i < v.NumField(); i++ {
			value := reflectValue.Field(i).Interface()
			tag := v.Field(i).Tag

			if tag != "" && tag != "-" {
				tagKey := tag.Get("json")
				jsonKey := strings.Split(tagKey, ",")[0]

				if jsonKey == field {
					mapped[field] = value
				}
			}
		}
	}

	return mapped
}
