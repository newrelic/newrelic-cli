package pipe

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStdinPipeReader(t *testing.T) {
	cases := map[string]struct {
		Input    string
		Expected string
	}{
		"singleInputReturnsSingleValue": {
			Input: `{
				"id": 1,
				"name": "Foo",
				"price": 123,
				"tags": [
					"Bar",
					"Eek"
				],
				"stock": {
					"warehouse": 300,
					"retail": 20
				}
			}`,
			Expected: `{ "id": 1, "name": "Foo", "price": 123, "tags": [ "Bar", "Eek" ], "stock": { "warehouse": 300, "retail": 20 } }`,
		},
		"arrayInputReturnsArrayValue": {
			Input: `[ 
				{ 
					"id": 1, 
					"name": "Foo", 
					"price": 123, 
					"tags": [ "Bar", "Eek" ], 
					"stock": { 
						"warehouse": 300, 
						"retail": 20
					}
				},
				{ 
					"id": 2, 
					"name": "Bar", 
					"price": 456, 
					"tags": [ "Oop", "Aah" ], 
					"stock": { 
						"warehouse": 450, 
						"retail": 50
					}
				},
				{ 
					"id": 3, 
					"name": "Baz", 
					"price": 789, 
					"tags": [ "Syn", "Ack" ], 
					"stock": { 
						"warehouse": 100, 
						"retail": 75
					}
				}
			]`,
			Expected: `[ { "id": 1, "name": "Foo", "price": 123, "tags": [ "Bar", "Eek" ], "stock": { "warehouse": 300, "retail": 20 } }, { "id": 2, "name": "Bar", "price": 456, "tags": [ "Oop", "Aah" ], "stock": { "warehouse": 450, "retail": 50 } }, { "id": 3, "name": "Baz", "price": 789, "tags": [ "Syn", "Ack" ], "stock": { "warehouse": 100, "retail": 75 } } ]`,
		},
	}

	for _, c := range cases {
		mockStdin := strings.NewReader(c.Input)
		reader := stdinPipeReader{input: mockStdin}
		result, err := reader.ReadPipe()

		assert.Equal(t, nil, err)
		assert.Equal(t, c.Expected, result)
	}

}

func TestJsonToFilteredMap(t *testing.T) {
	cases := map[string]struct {
		Input       string
		ExpectedErr error
		Expected    []map[string]string
	}{
		"noInputReturnsInvalidJsonError": {
			Input:       ``,
			Expected:    nil,
			ExpectedErr: errors.New("invalid JSON received by stdin"),
		},
		"singleValueReturnsCorrectOutput": {
			Input: `{ 
				"id": 1, 
				"name": "Foo", 
				"price": 123, 
				"tags": [ "Bar", "Eek" ], 
				"stock": { 
					"warehouse": 300, 
					"retail": 20
				}
			}`,
			Expected: []map[string]string{
				{
					"id":           "1",
					"stock.retail": "20",
				},
			},
			ExpectedErr: nil,
		},
		"arrayValueReturnsCorrectOutput": {
			Input: `[ 
				{ 
					"id": 1, 
					"name": "Foo", 
					"price": 123, 
					"tags": [ "Bar", "Eek" ], 
					"stock": { 
						"warehouse": 300, 
						"retail": 20
					}
				},
				{ 
					"id": 2, 
					"name": "Bar", 
					"price": 456, 
					"tags": [ "Oop", "Aah" ], 
					"stock": { 
						"warehouse": 450, 
						"retail": 50
					}
				},
				{ 
					"id": 3, 
					"name": "Baz", 
					"price": 789, 
					"tags": [ "Syn", "Ack" ], 
					"stock": { 
						"warehouse": 100, 
						"retail": 75
					}
				}
			]`,
			Expected: []map[string]string{
				{
					"id":           "1",
					"stock.retail": "20",
				},
				{
					"id":           "2",
					"stock.retail": "50",
				},
				{
					"id":           "3",
					"stock.retail": "75",
				},
			},
			ExpectedErr: nil,
		},
	}
	selectors := []string{"id", "stock.retail"}

	for _, c := range cases {
		results, err := jsonToFilteredMap(c.Input, selectors)

		assert.Equal(t, c.ExpectedErr, err)

		assert.Equal(t, c.Expected, results)
	}
}

func TestReadStdin(t *testing.T) {
	cases := map[string]struct {
		Input       string
		Expected    []map[string]string
		ExpectedErr error
	}{
		"invaildJsonThrowsError": {
			Input: `{ 
				broken = bad
			}`,
			Expected:    nil,
			ExpectedErr: errors.New("invalid JSON received by stdin"),
		},
		"singlarInputReturnsCorrectValue": {
			Input: `{ 
				"id": 1, 
				"name": "Foo", 
				"price": 123, 
				"tags": [ "Bar", "Eek" ], 
				"stock": { 
					"warehouse": 300, 
					"retail": 20
				}
			}`,
			Expected: []map[string]string{
				{
					"id":           "1",
					"stock.retail": "20",
				},
			},
			ExpectedErr: nil,
		},
		"arrayInputReturnsCorrectValue": {
			Input: `[ 
				{ 
					"id": 1, 
					"name": "Foo", 
					"price": 123, 
					"tags": [ "Bar", "Eek" ], 
					"stock": { 
						"warehouse": 300, 
						"retail": 20
					}
				},
				{ 
					"id": 2, 
					"name": "Bar", 
					"price": 456, 
					"tags": [ "Oop", "Aah" ], 
					"stock": { 
						"warehouse": 450, 
						"retail": 50
					}
				},
				{ 
					"id": 3, 
					"name": "Baz", 
					"price": 789, 
					"tags": [ "Syn", "Ack" ], 
					"stock": { 
						"warehouse": 100, 
						"retail": 75
					}
				}
			]`,
			Expected: []map[string]string{
				{
					"id":           "1",
					"stock.retail": "20",
				},
				{
					"id":           "2",
					"stock.retail": "50",
				},
				{
					"id":           "3",
					"stock.retail": "75",
				},
			},
			ExpectedErr: nil,
		},
	}
	selectors := []string{"id", "stock.retail"}

	for _, c := range cases {
		mockStdin := strings.NewReader(c.Input)

		results, err := readStdin(stdinPipeReader{input: mockStdin}, selectors)

		assert.Equal(t, c.ExpectedErr, err)

		assert.Equal(t, c.Expected, results)

	}
}

func TestGetPipeInputInnerFunc(t *testing.T) {
	cases := map[string]struct {
		Input    string
		Exists   bool
		Expected map[string][]string
	}{
		"noInputButStdinExists": {
			Input:    ``,
			Exists:   true,
			Expected: map[string][]string{},
		},
		"invalidJsonFromStdin": {
			Input: `{ 
				broken = bad
			}`,
			Exists:   true,
			Expected: map[string][]string{},
		},
		"inputGivenButNoStdinDetected": {
			Input: `{ 
				"id": 1, 
				"name": "Foo", 
				"price": 123, 
				"tags": [ "Bar", "Eek" ], 
				"stock": { 
					"warehouse": 300, 
					"retail": 20
				}
			}`,
			Exists:   false,
			Expected: map[string][]string{},
		},
		"singlarInputReturnsCorrectResults": {
			Input: `{ 
				"id": 1, 
				"name": "Foo", 
				"price": 123, 
				"tags": [ "Bar", "Eek" ], 
				"stock": { 
					"warehouse": 300, 
					"retail": 20
				}
			}`,
			Exists: true,
			Expected: map[string][]string{
				"id":           {"1"},
				"stock.retail": {"20"},
			},
		},
		"arrayInputReturnsCorrectResults": {
			Input: `[ 
				{ 
					"id": 1, 
					"name": "Foo", 
					"price": 123, 
					"tags": [ "Bar", "Eek" ], 
					"stock": { 
						"warehouse": 300, 
						"retail": 20
					}
				},
				{ 
					"id": 2, 
					"name": "Bar", 
					"price": 456, 
					"tags": [ "Oop", "Aah" ], 
					"stock": { 
						"warehouse": 450, 
						"retail": 50
					}
				},
				{ 
					"id": 3, 
					"name": "Baz", 
					"price": 789, 
					"tags": [ "Syn", "Ack" ], 
					"stock": { 
						"warehouse": 100, 
						"retail": 75
					}
				}
			]`,
			Exists: true,
			Expected: map[string][]string{
				"id":           {"1", "2", "3"},
				"stock.retail": {"20", "50", "75"},
			},
		},
	}
	selectors := []string{"id", "stock.retail"}

	for _, c := range cases {
		mockStdin := strings.NewReader(c.Input)

		results := getPipeInputInnerFunc(stdinPipeReader{input: mockStdin}, c.Exists, selectors)

		assert.Equal(t, c.Expected, results)
	}
}

func TestGetPipeInput(t *testing.T) {
	cases := map[string]struct {
		Input    string
		RunTwice bool
		Expected map[string][]string
	}{
		"setsInputCorrectlyIfRanOnce": {
			Input: `{ 
				"id": 1, 
				"name": "Foo", 
				"price": 123, 
				"tags": [ "Bar", "Eek" ], 
				"stock": { 
					"warehouse": 300, 
					"retail": 20
				}
			}`,
			RunTwice: false,
			Expected: map[string][]string{
				"id":           {"1"},
				"stock.retail": {"20"},
			},
		},
		"setsInputCorrectlyIfRanTwice": {
			Input: `{ 
				"id": 3, 
				"name": "Foo", 
				"price": 123, 
				"tags": [ "Bar", "Eek" ], 
				"stock": { 
					"warehouse": 300, 
					"retail": 30
				}
			}`,
			RunTwice: true,
			Expected: map[string][]string{
				"id":           {"3"},
				"stock.retail": {"30"},
			},
		},
	}

	selectors := []string{"id", "stock.retail"}

	for _, c := range cases {
		pipeInput = nil
		mockStdin := strings.NewReader(c.Input)

		getInput := getPipeInputFactory(stdinPipeReader{input: mockStdin}, func() bool { return true })

		getInput(selectors)

		assert.Equal(t, c.Expected, pipeInput)

		if c.RunTwice {
			getInput(selectors)
			assert.Equal(t, c.Expected, pipeInput)
		}
	}
}

func TestGet(t *testing.T) {
	cases := map[string]struct {
		Input         map[string][]string
		ExpectedOk    bool
		ExpectedValue []string
	}{
		"returnsNotOkIfPipeInputNotSet": {
			ExpectedValue: nil,
			ExpectedOk:    false,
			Input:         nil,
		},
		"returnsNotOkIfKeyNotPresent": {
			ExpectedValue: nil,
			ExpectedOk:    false,
			Input:         map[string][]string{},
		},
		"setsInputCorrectlyIfRanTwice": {
			ExpectedValue: []string{"3"},
			ExpectedOk:    true,
			Input: map[string][]string{
				"id":           {"3"},
				"stock.retail": {"30"},
			},
		},
	}

	for _, c := range cases {
		pipeInput = c.Input

		value, ok := Get("id")

		assert.Equal(t, c.ExpectedOk, ok)

		assert.Equal(t, c.ExpectedValue, value)
	}
}
