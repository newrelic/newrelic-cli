// +build integration

package extensions

import (
	"bufio"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExtensionStart(t *testing.T) {
	echo := &Manifest{
		Command: "echo",
	}

	e := New(
		echo,
		WithArgs("Hello world\n"),
	)

	stdout, err := e.cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}

	scanner := bufio.NewScanner(stdout)

	e.Start()

	scanner.Scan()
	assert.Equal(t, "Hello world", scanner.Text())
}

func TestExtensionTimeout(t *testing.T) {
	echo := &Manifest{
		Command: "sleep",
	}

	e := New(
		echo,
		WithTimeout(time.Duration(1*time.Second)),
		WithArgs("1"),
	)

	e.Start()

	<-e.ctx.Done()

	assert.IsType(t, context.DeadlineExceeded, e.ctx.Err())
}
