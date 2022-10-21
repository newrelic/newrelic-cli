package execution

import (
	"fmt"
	"strconv"
)

type OutputParser struct {
	output map[string]interface{}
}

func NewOutputParser(output map[string]interface{}) *OutputParser {
	return &OutputParser{
		output: output,
	}
}

func (op *OutputParser) EntityGUID() string {
	if val, ok := op.output["EntityGuid"]; ok {
		return fmt.Sprintf("%s", val)
	}
	return ""
}

func (op *OutputParser) Metadata() map[string]string {
	metadata, ok := op.output["Metadata"].(map[string]interface{})
	if ok {
		result := map[string]string{}
		for k := range metadata {
			if v, ok := metadata[k].(string); ok {
				result[k] = v
			}
		}
		return result
	}
	return nil
}

// Get indicator if recipe wants to log the recipe output
func (op *OutputParser) LogRecipeOutput() bool {
	metadata := op.Metadata()
	if val, ok := metadata["LogRecipeOutput"]; ok {
		if flag, err := strconv.ParseBool(val); err == nil {
			return flag
		}
	}
	return false
}

func (op *OutputParser) FailedRecipeOutput() string {
	if val, ok := op.output["FailedRecipeOutput"]; ok {
		return val.(string)
	}
	return ""
}
