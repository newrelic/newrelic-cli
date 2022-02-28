package install

import "github.com/newrelic/newrelic-cli/internal/install/recipes"

type BundlerBuilder struct {
	coreRecipes []*recipes.BundleRecipe
}

func NewBundlerBuilder() *BundlerBuilder {
	return &BundlerBuilder{}
}

func (bb *BundlerBuilder) WithCoreRecipe(name string) *BundlerBuilder {
	coreRecipes := []*recipes.BundleRecipe{
		{
			Recipe: recipes.NewRecipeBuilder().Name(name).Build(),
		},
	}
	bb.coreRecipes = coreRecipes
	return bb
}

func (bb *BundlerBuilder) WithAdditionalRecipe(name string) *BundlerBuilder {
	coreRecipes := []*recipes.BundleRecipe{
		{
			Recipe: recipes.NewRecipeBuilder().Name(name).Build(),
		},
	}
	bb.coreRecipes = coreRecipes
	return bb
}

func (bb *BundlerBuilder) Build() RecipeBundler {
	return &MockBundler{
		coreRecipes: bb.coreRecipes,
	}
}

type MockBundler struct {
	coreRecipes []*recipes.BundleRecipe
}

func (mb *MockBundler) CreateCoreBundle() *recipes.Bundle {
	bundle := &recipes.Bundle{
		BundleRecipes: mb.coreRecipes,
	}
	return bundle
}
func (mb *MockBundler) CreateAdditionalTargetedBundle(name []string, recipePaths []string) (*recipes.Bundle, error) {

	bundle := &recipes.Bundle{
		BundleRecipes: []*recipes.BundleRecipe{
			{
				Recipe: recipes.NewRecipeBuilder().ID("1").Name("Additional_Target_Recipe_1").Build(),
			},
		},
	}
	return bundle, nil
}

func (mb *MockBundler) CreateAdditionalGuidedBundle() *recipes.Bundle {

	bundle := &recipes.Bundle{
		BundleRecipes: []*recipes.BundleRecipe{
			{
				Recipe: recipes.NewRecipeBuilder().ID("2").Name("Additional_Guided_Recipe_1").Build(),
			},
		},
	}
	return bundle
}
