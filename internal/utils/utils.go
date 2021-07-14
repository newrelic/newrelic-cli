package utils

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/mitchellh/go-homedir"

	log "github.com/sirupsen/logrus"
)

var (
	SignalCtx = getSignalContext()
)

func getSignalContext() context.Context {
	ch := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-ch
		log.Debugf("signal received: %s", sig)
		cancel()
	}()
	return ctx
}

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

// MinOf returns the minimum int value provided.
func MinOf(vars ...int) int {
	min := vars[0]

	for _, i := range vars {
		if min > i {
			min = i
		}
	}

	return min
}

// GetTimestamp returns the current epoch timestamp in seconds.
func GetTimestamp() int64 {
	return time.Now().Unix()
}

// MakeRange generates a slice of sequential integers.
func MakeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

// Base64Encode base 64 encodes a string.
func Base64Encode(data string) string {
	return b64.StdEncoding.EncodeToString([]byte(data))
}

// Standard way to check for stdin in most environments (https://stackoverflow.com/questions/22563616/determine-if-stdin-has-data-with-go)
func StdinExists() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) == 0
}

func StringInSlice(str string, slice []string) bool {
	for _, s := range slice {
		if str == s {
			return true
		}
	}

	return false
}

func IntSliceToStringSlice(in []int) (out []string) {
	for _, i := range in {
		out = append(out, strconv.Itoa(i))
	}

	return out
}

// Obfuscate receives a string, and replaces everything after the first 8
// characters with an asterisk before returning the result.
func Obfuscate(input string) string {
	result := make([]string, len(input))
	parts := strings.Split(input, "")

	for i, x := range parts {
		if i < 8 {
			result[i] = x
		} else {
			result[i] = "*"
		}
	}

	return strings.Join(result, "")
}
