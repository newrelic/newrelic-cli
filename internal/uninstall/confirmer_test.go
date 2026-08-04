//go:build unit

package uninstall

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/ux"
)

func TestSurveyConfirmer_ReturnsPrompterAnswer(t *testing.T) {
	prompter := ux.NewMockPrompter()
	prompter.PromptYesNoVal = true
	c := &SurveyConfirmer{prompter: prompter}

	ok, err := c.Confirm("continue?")

	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, 1, prompter.PromptYesNoCallCount)
}

func TestSurveyConfirmer_PropagatesPrompterError(t *testing.T) {
	prompter := ux.NewMockPrompter()
	prompter.PromptYesNoErr = errors.New("boom")
	c := &SurveyConfirmer{prompter: prompter}

	_, err := c.Confirm("continue?")

	require.Error(t, err)
}

func TestNewSurveyConfirmer(t *testing.T) {
	require.NotNil(t, NewSurveyConfirmer())
}
