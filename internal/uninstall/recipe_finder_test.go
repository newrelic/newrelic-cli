//go:build unit

package uninstall

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestRepositoryRecipeFinder_FindsHostCompatibleRecipe(t *testing.T) {
	loader := func() ([]*types.OpenInstallationRecipe, error) {
		return []*types.OpenInstallationRecipe{
			{Name: "infra-agent"},
		}, nil
	}
	manifest := &types.DiscoveryManifest{OS: "linux"}

	f := NewRepositoryRecipeFinder(loader, manifest)
	r, err := f.FindRecipe(context.Background(), "infra-agent")

	require.NoError(t, err)
	require.Equal(t, "infra-agent", r.Name)
}

func TestRepositoryRecipeFinder_UnknownNameIsNotFound(t *testing.T) {
	loader := func() ([]*types.OpenInstallationRecipe, error) {
		return []*types.OpenInstallationRecipe{
			{Name: "infra-agent"},
		}, nil
	}
	manifest := &types.DiscoveryManifest{OS: "linux"}

	f := NewRepositoryRecipeFinder(loader, manifest)
	_, err := f.FindRecipe(context.Background(), "does-not-exist")

	require.True(t, errors.Is(err, ErrRecipeNotFound))
}

func TestRepositoryRecipeFinder_KnownButWrongHostIsUnsupported(t *testing.T) {
	loader := func() ([]*types.OpenInstallationRecipe, error) {
		return []*types.OpenInstallationRecipe{
			{
				Name: "windows-only-recipe",
				InstallTargets: []types.OpenInstallationRecipeInstallTarget{
					{Os: types.OpenInstallationOperatingSystemTypes.WINDOWS},
				},
			},
		}, nil
	}
	manifest := &types.DiscoveryManifest{OS: "linux"}

	f := NewRepositoryRecipeFinder(loader, manifest)
	_, err := f.FindRecipe(context.Background(), "windows-only-recipe")

	require.True(t, errors.Is(err, ErrRecipeUnsupportedOnHost))
}

func TestRepositoryRecipeFinder_LoaderErrorIsPropagated(t *testing.T) {
	loadErr := errors.New("boom")
	loader := func() ([]*types.OpenInstallationRecipe, error) {
		return nil, loadErr
	}
	manifest := &types.DiscoveryManifest{OS: "linux"}

	f := NewRepositoryRecipeFinder(loader, manifest)
	_, err := f.FindRecipe(context.Background(), "anything")

	require.ErrorIs(t, err, loadErr)
}
