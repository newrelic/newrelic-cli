package execution

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// odd hack that works
func cleanup() {
	fmt.Printf("\n")
}

func cleanupWithFile(filename string) {
	os.Remove(filename)
	cleanup()
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
	file, err := os.CreateTemp("", "some-test-file")
	assert.NoError(t, err)
	outputFile := file.Name()

	_, err = file.WriteString("error line one\n")
	assert.NoError(t, err)
	_, err = file.WriteString("error line two\n")
	assert.NoError(t, err)
	_, err = file.WriteString("error line three\n")
	assert.NoError(t, err)

	rlf := NewRecipeLogForwarder()
	rlf.SendLogsToNewRelic(outputFile, "test-recipe")

	assert.Equal(t, 3, len(rlf.LogEntries))
	cleanupWithFile(outputFile)
}
