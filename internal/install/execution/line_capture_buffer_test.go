package execution

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLineCaptureBuffer(t *testing.T) {
	w := bytes.NewBufferString("")
	b := NewLineCaptureBuffer(w)
	_, err := b.Write([]byte("abc\n123\ndef"))
	assert.NoError(t, err)

	require.Equal(t, "123", b.LastFullLine)
	require.Equal(t, "def", b.Current())
	require.Equal(t, "abc\n123\ndef", w.String())
}
