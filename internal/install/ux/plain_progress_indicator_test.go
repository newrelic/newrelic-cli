package ux

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPlainProgressIndicator_interface(t *testing.T) {
	var r ProgressIndicator = NewPlainProgress()
	require.NotNil(t, r)
}
