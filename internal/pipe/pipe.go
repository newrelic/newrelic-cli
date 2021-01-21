// Package pipe provides a simple API to read and retrieve values
// from stdin to use in Cobra commands. Public API consists of
// GetInput, which reads stdin, Exists, which checks for value
// existence, and Get for retrieving existing values.
package pipe

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

// Created Interface and struct to surround io.Reader for easy mocking
type pipeReader interface {
	ReadPipe() (string, error)
}

type stdinPipeReader struct {
	input io.Reader
}

func (spr stdinPipeReader) ReadPipe() (string, error) {
	text := ""
	var err error = nil
	scanner := bufio.NewScanner(spr.input)

	for scanner.Scan() {
		text = fmt.Sprintf("%s%s ", text, strings.TrimSpace(scanner.Text()))
	}

	if scanErr := scanner.Err(); scanErr != nil {
		err = scanErr
	}

	return strings.TrimSpace(text), err
}

func jsonToFilteredMap(r string, selectors []string) ([]map[string]string, error) {
	text := r

	if !gjson.Valid(text) {
		return nil, errors.New("invalid JSON received by stdin")
	}

	// always start with an array of values
	if strings.HasPrefix(text, "{") {
		text = fmt.Sprintf("[ %s ]", text)
	}

	// returns []gjson.Result
	jsonArray := gjson.Parse(text).Array()

	resultsArray := make([]map[string]string, len(jsonArray))

	for index, resultObj := range jsonArray {
		resultMap := make(map[string]string)
		for _, selector := range selectors {
			if value := resultObj.Get(selector); value.Exists() {
				// Convert every value to a string
				resultMap[selector] = value.String()
			}
		}
		resultsArray[index] = resultMap
	}

	return resultsArray, nil
}

func readStdin(pipe pipeReader, selectorList []string) ([]map[string]string, error) {
	jsonString, pipeErr := pipe.ReadPipe()
	if pipeErr != nil {
		return nil, pipeErr
	}

	filteredMap, mapErr := jsonToFilteredMap(jsonString, selectorList)
	if mapErr != nil {
		return nil, mapErr
	}

	return filteredMap, nil
}

// Standard way to check for stdin in most environments (https://stackoverflow.com/questions/22563616/determine-if-stdin-has-data-with-go)
func pipeInputExists() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) == 0
}

// getPipeInputInnerFunc is the function returned with arguments by
// getPipeInputFactory that becomes GetInput. Private to package, main
// function to be tested.
func getPipeInputInnerFunc(pipe pipeReader, pipeInputExists bool, acceptedPipeInput []string) map[string][]string {
	if pipeInputExists {
		pipeInputMap := map[string][]string{}
		inputArray, err := readStdin(pipe, acceptedPipeInput)
		if err != nil {
			log.Error(err)
			return map[string][]string{}
		}
		for _, key := range acceptedPipeInput {
			var collectedItemsForKey []string
			for _, value := range inputArray {
				collectedItemsForKey = append(collectedItemsForKey, value[key])
			}
			pipeInputMap[key] = collectedItemsForKey
		}
		return pipeInputMap
	}
	return map[string][]string{}
}

// getPipeInputFactory is a factory method to create GetInput. Allows for
// extensive testing options using dependency injection. Only ever called
// once on the first import of the package. Private to the package.
func getPipeInputFactory(pipe pipeReader, predicate func() bool) func([]string) {
	return func(acceptedPipeInput []string) {
		if pipeInput == nil {
			pipeInput = getPipeInputInnerFunc(pipe, predicate(), acceptedPipeInput)
		}
	}
}

var pipeInput map[string][]string

// GetInput takes a slice of gjson selectors (https://github.com/tidwall/gjson/blob/master/SYNTAX.md)
// as an argument. When ran once at the top the init function, GetInput
// stores those desired json values from stdin. The existence of and values
// of those stdin json keys can then be retrieved using the public Exists and
// Get methods, respectively.
var GetInput = getPipeInputFactory(stdinPipeReader{input: os.Stdin}, pipeInputExists)

// Get is the only API provided to retrieve values from stdin json. Get
// is designed to be used in the cobra command itself, when any required
// value fed into stdin is needed. If stdin is empty or the value specified
// does not exist, Get will return nil for the value and false for the ok
// check.
func Get(inputKey string) ([]string, bool) {
	if pipeInput == nil {
		return nil, false
	}
	value, ok := pipeInput[inputKey]
	return value, ok
}

// Exists is meant to be used in the init function, after GetInput is called,
// to determine if a required Cobra flag needs to be declared. If the inputKey
// exists in the Pipe module, then you can skip the required flag declaration
// and used the piped in values instead.
func Exists(inputKey string) bool {
	_, ok := Get(inputKey)
	return ok
}
