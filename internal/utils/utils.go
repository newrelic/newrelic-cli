package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mitchellh/go-homedir"

	log "github.com/sirupsen/logrus"
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

// LogIfError wraps the err nil check to cleanup the code.
// Logs at Error level
func LogIfError(err error) {
	if err != nil {
		log.Error(err)
	}
}

// LogIfFatal wraps the err nil check to cleanup the code.
// Logs at Fatal level
func LogIfFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// GetDefaultConfigDirectory returns the full path to the .newrelic
// directory within the user's home directory.
func GetDefaultConfigDirectory() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/.newrelic", home), nil
}
