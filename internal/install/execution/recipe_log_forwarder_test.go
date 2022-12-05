package execution

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// odd hack that works
func cleanup() {
	fmt.Printf("\n")
}

func TestSendLogsIsTrueWhenyYOrNothingEntered(t *testing.T) {
	rlf := NewRecipeLogForwarder()
	var buffer bytes.Buffer

	buffer.WriteString("Y")
	assert.Equal(t, true, rlf.PromptUserToSendLogs(&buffer))

	buffer.WriteString("y")
	assert.True(t, rlf.PromptUserToSendLogs(&buffer))

	buffer.WriteString(" ")
	assert.True(t, rlf.PromptUserToSendLogs(&buffer))
	cleanup()
}

func TestSendLogsIsFalseWhenAnythingButyYandNothingEntered(t *testing.T) {
	rlf := NewRecipeLogForwarder()
	var buffer bytes.Buffer

	buffer.WriteString("n")
	assert.False(t, rlf.PromptUserToSendLogs(&buffer))

	buffer.WriteString("N")
	assert.False(t, rlf.PromptUserToSendLogs(&buffer))

	buffer.WriteString("omg")
	assert.False(t, rlf.PromptUserToSendLogs(&buffer))
	cleanup()
}

func TestSendLogsToNewRelicBuildsLogEntryForEachLogLine(t *testing.T) {
	rlf := NewRecipeLogForwarder()
	rlf.SendLogsToNewRelic("test-recipe", []string{"error line one\n", "error line two\n", "error line three\n"})

	assert.Equal(t, 3, len(rlf.LogEntries))
}
