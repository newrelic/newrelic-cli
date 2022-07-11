package execution

import "fmt"

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
