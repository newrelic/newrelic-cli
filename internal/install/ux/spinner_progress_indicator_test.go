package ux

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/config"
)

func TestConnectToPlatformShouldSuccess(t *testing.T) {
	msg := "msg"
	pi := NewSpinnerProgressIndicator()

	stdOut := captureLoggingOutput(func() {
		pi.Start(msg)
	})
	assert.True(t, strings.Contains(stdOut, msg))
}

func captureLoggingOutput(f func()) string {
	var buf bytes.Buffer
	existingLogger := config.Logger
	existingLogger.SetOutput(&buf)
	existingLogger.SetLevel(logrus.DebugLevel)
	f()
	existingLogger.SetOutput(os.Stderr)
	return buf.String()
}
