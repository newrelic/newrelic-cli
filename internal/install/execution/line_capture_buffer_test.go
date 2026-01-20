package execution

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLineCaptureBuffer(t *testing.T) {
	tests := []struct {
		name           string
		inputs         []string // Multiple writes if needed
		expectedLast   string   // Expected LastFullLine
		expectedCurr   string   // Expected Current()
		expectedOutput string   // Expected final buffer content
	}{
		{
			name:           "Basic input with multiple lines",
			inputs:         []string{"abc\n123\ndef"},
			expectedLast:   "123",
			expectedCurr:   "def",
			expectedOutput: "abc\n123\ndef",
		},
		{
			name:           "Input with trailing empty line",
			inputs:         []string{"abc\n123\ndef\n \n"},
			expectedLast:   "def",
			expectedCurr:   "",
			expectedOutput: "abc\n123\ndef\n \n",
		},
		{
			name:           "Multiple write operations",
			inputs:         []string{"abc\n", "123\n", "def\n", "nope"},
			expectedLast:   "def",
			expectedCurr:   "nope",
			expectedOutput: "abc\n123\ndef\nnope",
		},
		{
			name:           "Empty input",
			inputs:         []string{""},
			expectedLast:   "",
			expectedCurr:   "",
			expectedOutput: "",
		},
		{
			name:           "Single line without newline",
			inputs:         []string{"single"},
			expectedLast:   "",
			expectedCurr:   "single",
			expectedOutput: "single",
		},
		{
			name:           "Multiple empty lines",
			inputs:         []string{"\n\n\n"},
			expectedLast:   "",
			expectedCurr:   "",
			expectedOutput: "\n\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := bytes.NewBufferString("")
			b := NewLineCaptureBuffer(w)

			// Process all inputs
			for _, input := range tt.inputs {
				_, err := b.Write([]byte(input))
				assert.NoError(t, err, "Write operation should not fail")
			}

			// Verify the results
			require.Equal(t, tt.expectedLast, b.LastFullLine, "LastFullLine mismatch")
			require.Equal(t, tt.expectedCurr, b.Current(), "Current() mismatch")
			require.Equal(t, tt.expectedOutput, w.String(), "Buffer content mismatch")
		})
	}
}
