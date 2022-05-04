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
