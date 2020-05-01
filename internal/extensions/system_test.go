// +build integration

package extensions

import (
	"bufio"
	"context"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExtensionStart(t *testing.T) {
	echo := &Manifest{
		Command: "echo",
	}

	e, err := New(
		echo,
		WithArgs("Hello world\n"),
	)
	if err != nil {
		t.Fatal(err)
	}

	scanner := bufio.NewScanner(e.StdoutPipe)

	if err := e.Start(); err != nil {
		t.Fatal(err)
	}

	scanner.Scan()
	assert.Equal(t, "Hello world", scanner.Text())
}

func TestExtensionCleanShutdown(t *testing.T) {
	echo := &Manifest{
		Command: "echo",
	}

	e, err := New(echo)
	if err != nil {
		t.Fatal(err)
	}

	if err := e.Start(); err != nil {
		t.Fatal(err)
	}

	<-e.DoneChan
	assert.Nil(t, e.Err())
}

func TestExtensionError(t *testing.T) {
	cmd := &Manifest{
		Command: "false",
	}

	e, err := New(cmd)
	if err != nil {
		t.Fatal(err)
	}

	if err = e.Start(); err != nil {
		t.Fatal(err)
	}

	<-e.DoneChan
	assert.NotNil(t, e.Err())
	assert.IsType(t, &exec.ExitError{}, e.Err())
}

func TestExtensionTimeout(t *testing.T) {
	echo := &Manifest{
		Command: "sleep",
	}

	e, err := New(
		echo,
		WithTimeout(time.Duration(100*time.Millisecond)),
		WithArgs("1"),
	)
	if err != nil {
		t.Fatal(err)
	}

	if err = e.Start(); err != nil {
		t.Fatal(err)
	}

	<-e.DoneChan
	assert.NotNil(t, e.Err())
	assert.IsType(t, e.Err(), context.DeadlineExceeded)
}

func TestExtensionCancel(t *testing.T) {
	echo := &Manifest{
		Command: "sleep",
	}

	e, err := New(
		echo,
		WithTimeout(time.Duration(100*time.Millisecond)),
		WithArgs("infinity"),
	)
	if err != nil {
		t.Fatal(err)
	}

	if err = e.Start(); err != nil {
		t.Fatal(err)
	}

	e.CancelFunc()

	<-e.DoneChan
}
