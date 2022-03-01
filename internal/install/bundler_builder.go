package install

import "github.com/newrelic/newrelic-cli/internal/install/recipes"

type BundlerBuilder struct {
	coreRecipes       []*recipes.BundleRecipe
	additionalRecipes []*recipes.BundleRecipe
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
	additionalRecipes := []*recipes.BundleRecipe{
		{
			Recipe: recipes.NewRecipeBuilder().Name(name).Build(),
		},
	}
	bb.additionalRecipes = additionalRecipes
	return bb
}

func (bb *BundlerBuilder) Build() RecipeBundler {
	return &MockBundler{
		coreRecipes:       bb.coreRecipes,
		additionalRecipes: bb.additionalRecipes,
	}
}

type MockBundler struct {
	coreRecipes       []*recipes.BundleRecipe
	additionalRecipes []*recipes.BundleRecipe
}

func (mb *MockBundler) CreateCoreBundle() *recipes.Bundle {
	bundle := &recipes.Bundle{
		BundleRecipes: mb.coreRecipes,
	}
	return bundle
}
func (mb *MockBundler) CreateAdditionalTargetedBundle(names []string) *recipes.Bundle {

	bundle := &recipes.Bundle{
		//		BundleRecipes: mb.additionalRecipes,
	}

	for _, r := range mb.additionalRecipes {
		for _, n := range names {
			if r.Recipe.Name == n {
				bundle.AddRecipe(r)
			}
		}
	}

	return bundle
}

func (mb *MockBundler) CreateAdditionalGuidedBundle() *recipes.Bundle {

	bundle := &recipes.Bundle{
		BundleRecipes: mb.additionalRecipes,
	}
	return bundle
}
