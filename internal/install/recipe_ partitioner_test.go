package install

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

var withRecipes = []types.OpenInstallationRecipe{
	{
		Name:           "testRecipe",
		ValidationNRQL: "testNrql",
		Dependencies:   []string{types.InfraAgentRecipeName},
	},
	{
		Name:           types.InfraAgentRecipeName,
		ValidationNRQL: "testNrql",
	},
	{
		Name: types.LoggingRecipeName,
	},
	{
		Name: types.GoldenRecipeName,
	},
}

var withoutRecipes = []types.OpenInstallationRecipe{}

const expectedPartitionCount = 2

func TestRecipePartition_Any(t *testing.T) {

	tests := []struct {
		name     string
		recipes  []types.OpenInstallationRecipe
		excepted bool
	}{
		{"With recipes", withRecipes, true},
		{"Without recipes", withoutRecipes, false},
	}

	sut := recipePartition{}

	for _, e := range tests {
		t.Run(e.name, func(t *testing.T) {
			sut.recipes = e.recipes
			actual := sut.any()
			if actual != e.excepted {
				t.Errorf("actual %v excepted %v", actual, e.excepted)
			}
		})
	}
}

func TestRecipePartition_Partition(t *testing.T) {

	sut := coreRecipePartition
	other := sut.partition(withRecipes)

	if len(sut.recipes) != 3 {
		t.Errorf("Partition core actual %d expected %d", len(sut.recipes), 1)
	}

	otherCount := len(withRecipes) - len(sut.recipes)

	if otherCount != len(other) {
		t.Errorf("Partition other actual %d expected %d", len(other), otherCount)
	}

}

func TestRecipePartitions_New(t *testing.T) {

	sut := *newRecipePartitions(withRecipes)
	if len(sut) != expectedPartitionCount {
		t.Errorf("Partiton count excepted %d actual %d", expectedPartitionCount, len(sut))
	}

	for _, p := range sut {
		if p.name == coreRecipePartition.name {
			if len(p.recipes) != 3 {
				t.Errorf("Partition Core Recipe count excepted %d actual %d", 1, len(p.recipes))
			}
		} else if p.name == otherRecipePartition.name {
			if len(p.recipes) != 1 {
				t.Errorf("Partition Other Recipe count excepted %d actual %d", 1, len(p.recipes))
			}
		}
	}
}
