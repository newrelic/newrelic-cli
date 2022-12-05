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

func (op *OutputParser) IsCapturedCliOutput() bool {
	capturedOutput := false
	data, ok := op.metadataWithKey("CapturedCliOutput")
	if ok {
		valAsBool, err := strconv.ParseBool(data)
		if err == nil {
			capturedOutput = valAsBool
		}
	}

	return capturedOutput
}

func (op *OutputParser) metadataWithKey(key string) (string, bool) {
	metadata, ok := op.output["Metadata"].(map[string]interface{})
	if ok {
		data, exists := metadata[key]
		if exists {
			return data.(string), ok
		}
	}
	return "", ok
}

func (op *OutputParser) AddMetadata(key string, value string) {
	output := op.output
	if output == nil {
		op.output = make(map[string]interface{})
	}

	metadata, ok := op.output["Metadata"].(map[string]interface{})
	if !ok {
		metadata = make(map[string]interface{})
	}

	metadata[key] = value
	op.output["Metadata"] = metadata
}
