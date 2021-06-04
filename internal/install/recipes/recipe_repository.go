package recipes

import (
	"math"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type RecipeRepository struct {
	RecipeLoaderFunc func() ([]types.OpenInstallationRecipe, error)
	recipes          []types.OpenInstallationRecipe
}

type recipeMatch struct {
	matchCount int
	recipe     types.OpenInstallationRecipe
}

// NewRecipeRepository returns a new instance of types.RecipeRepository.
func NewRecipeRepository(loaderFunc func() ([]types.OpenInstallationRecipe, error)) *RecipeRepository {
	rr := RecipeRepository{
		RecipeLoaderFunc: loaderFunc,
		recipes:          nil,
	}

	return &rr
}

func (rf *RecipeRepository) FindAll(m types.DiscoveryManifest) ([]types.OpenInstallationRecipe, error) {
	results := []types.OpenInstallationRecipe{}
	matchRecipes := make(map[string][]recipeMatch)
	hostMap := getHostMap(m)

	log.Debugf("Find all recipes available for host %+v", hostMap)

	if rf.recipes == nil {
		recipes, err := rf.RecipeLoaderFunc()
		if err != nil {
			return nil, err
		}

		rf.recipes = recipes
	}

	for _, recipe := range rf.recipes {
		matchTargetCount := []int{}

		for _, rit := range recipe.InstallTargets {
			matchCount := 0
			for k, v := range getRecipeTargetMap(rit) {
				isValueMatching := matchRecipeCriteria(hostMap, k, v)
				if isValueMatching {
					matchCount++
				} else {
					log.Debugf("recipe %s defines %s but input did not provide a match", recipe.Name, k)
					matchCount = 0
					break
				}
			}
			if matchCount > 0 {
				matchTargetCount = append(matchTargetCount, matchCount)
			}
		}

		if len(recipe.InstallTargets) == 0 || len(matchTargetCount) > 0 {
			maxMatchTargetCount := 0
			if len(matchTargetCount) > 0 {
				maxMatchTargetCount = mathMax(matchTargetCount)
			}

			match := recipeMatch{
				recipe:     recipe,
				matchCount: maxMatchTargetCount,
			}
			if _, ok := matchRecipes[recipe.Name]; !ok {
				matches := []recipeMatch{match}
				matchRecipes[recipe.Name] = matches
			} else {
				matchRecipes[recipe.Name] = append(matchRecipes[recipe.Name], match)
			}
		}
	}

	for _, matches := range matchRecipes {
		if len(matches) > 0 {
			match := findMaxMatch(matches)
			singleRecipe := match.recipe
			results = append(results, singleRecipe)
		}
	}

	return results, nil
}

func findMaxMatch(matches []recipeMatch) recipeMatch {
	var result *recipeMatch

	for _, match := range matches {
		if result == nil {
			result = &match
		} else {
			if match.matchCount > result.matchCount {
				result = &match
			}
		}
	}

	return *result
}

func mathMax(numbers []int) int {
	result := math.MinInt32
	for _, number := range numbers {
		if number > result {
			result = number
		}
	}
	return result
}

func matchRecipeCriteria(hostMap map[string]string, rkey string, rvalue string) bool {
	if rvalue == "" {
		return true
	}

	if val, ok := hostMap[rkey]; ok {
		return strings.EqualFold(val, rvalue)
	}

	return false
}

func getHostMap(m types.DiscoveryManifest) map[string]string {
	hostMap := map[string]string{
		"KernelArch":      m.KernelArch,
		"KernelVersion":   m.KernelVersion,
		"OS":              m.OS,
		"Platform":        m.Platform,
		"PlatformFamily":  m.PlatformFamily,
		"PlatformVersion": m.PlatformVersion,
	}
	return hostMap
}

func getRecipeTargetMap(rit types.OpenInstallationRecipeInstallTarget) map[string]string {
	targetMap := map[string]string{
		"KernelArch":      rit.KernelArch,
		"KernelVersion":   rit.KernelVersion,
		"OS":              string(rit.Os),
		"Platform":        string(rit.Platform),
		"PlatformFamily":  string(rit.PlatformFamily),
		"PlatformVersion": rit.PlatformVersion,
	}
	return targetMap
}
