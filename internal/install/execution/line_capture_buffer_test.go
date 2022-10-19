package execution

import (
	"bytes"
	"io/ioutil"
	"os"
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

func TestLineCaptureBufferToFile(t *testing.T) {

	w := bytes.NewBufferString("")
	outputFile, _ := ioutil.TempFile("", "some-test-file_")
	defer outputFile.Close()
	defer os.Remove(outputFile.Name())

	b := NewLineCaptureToFileBuffer(w, outputFile)
	_, err := b.Write([]byte("abc\n123\ndef"))
	assert.NoError(t, err)
	require.Equal(t, "123", b.LastFullLine)
	require.Equal(t, "def", b.Current())
	require.Equal(t, "abc\n123\ndef", w.String())

	data, err := os.ReadFile(outputFile.Name())
	assert.NoError(t, err)
	require.Equal(t, "abc\n123\n", string(data)) // only writes full lines to output file

}
