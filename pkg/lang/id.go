package lang

import "fmt"

// Language identifiers.
const (
	Unknown ID = "unknown"
	Java    ID = "java"
	Go      ID = "go"
)

// ID language unique identifier.
type ID string

// IntegrationName generates the expected language specific introspection integration name.
func (i ID) IntegrationName() string {
	return fmt.Sprintf("nri-lsi-%s", i)
}
