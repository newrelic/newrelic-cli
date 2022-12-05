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

func TestLineCaptureBufferCapturesEntireOutput(t *testing.T) {
	w := bytes.NewBufferString("")

	b := NewLineCaptureBuffer(w)
	_, err := b.Write([]byte("abc\n"))
	assert.NoError(t, err)

	_, err = b.Write([]byte("123\n"))
	assert.NoError(t, err)

	_, err = b.Write([]byte("def\n"))
	assert.NoError(t, err)

	_, err = b.Write([]byte("nope"))
	assert.NoError(t, err)

	require.Equal(t, "def", b.LastFullLine)
	require.Equal(t, "nope", b.Current())

	require.Equal(t, len(b.fullRecipeOutput), 3)

}
