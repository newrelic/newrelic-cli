package recipes

import (
	"math"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

var (
	kernelArch      string = "KernelArch"
	kernelVersion   string = "KernelVersion"
	oS              string = "OS"
	platform        string = "Platform"
	platformFamily  string = "PlatformFamily"
	platformVersion string = "PlatformVersion"
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

	if rf.recipes == nil {
		recipes, err := rf.RecipeLoaderFunc()
		if err != nil {
			return nil, err
		}
		log.Debugf("Loaded %d recipes", len(recipes))

		rf.recipes = recipes
	}

	log.Debugf("Find all available out of %d recipes for host %+v", len(rf.recipes), hostMap)
	log.Debugf("All recipes: %+v", rf.recipes)

	for _, recipe := range rf.recipes {
		matchTargetCount := []int{}

		for _, rit := range recipe.InstallTargets {
			matchCount := 0
			for k, v := range getRecipeTargetMap(rit) {
				if v == "" {
					continue
				}
				isValueMatching := matchRecipeCriteria(hostMap, k, v)
				if isValueMatching {
					log.Debugf("matching recipe %s field name %s and value %s using hostMap %+v", recipe.Name, k, v, hostMap)
					matchCount++
				} else {
					log.Debugf("recipe %s defines %s=%s but input did not provide a match using hostMap %+v", recipe.Name, k, v, hostMap)
					matchCount = -1
					break
				}
			}
			if matchCount >= 0 {
				matchTargetCount = append(matchTargetCount, matchCount)
			}
		}

		if len(recipe.InstallTargets) == 0 || len(matchTargetCount) > 0 {
			maxMatchTargetCount := 0
			if len(matchTargetCount) > 0 {
				maxMatchTargetCount = mathMax(matchTargetCount)
			}
			log.Debugf("Recipe InstallTargetsCount %d and maxMatchCount %d", len(recipe.InstallTargets), maxMatchTargetCount)

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
			log.Debugf("Add result for recipe name %s with targets %+v", match.recipe, match.recipe.InstallTargets)
			results = append(results, singleRecipe)
		}
	}

	return results, nil
}

func findMaxMatch(matches []recipeMatch) recipeMatch {
	var result *recipeMatch

	for _, match := range matches {
		if result == nil {
			result = &recipeMatch{
				recipe:     match.recipe,
				matchCount: match.matchCount,
			}
		} else {
			if match.matchCount > result.matchCount {
				result = &recipeMatch{
					recipe:     match.recipe,
					matchCount: match.matchCount,
				}
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
	if val, ok := hostMap[rkey]; ok {
		return strings.EqualFold(val, rvalue)
	}

	return false
}

func getHostMap(m types.DiscoveryManifest) map[string]string {
	hostMap := map[string]string{
		kernelArch:      m.KernelArch,
		kernelVersion:   m.KernelVersion,
		oS:              m.OS,
		platform:        m.Platform,
		platformFamily:  m.PlatformFamily,
		platformVersion: m.PlatformVersion,
	}
	return hostMap
}

func getRecipeTargetMap(rit types.OpenInstallationRecipeInstallTarget) map[string]string {
	targetMap := map[string]string{
		kernelArch:      rit.KernelArch,
		kernelVersion:   rit.KernelVersion,
		oS:              string(rit.Os),
		platform:        string(rit.Platform),
		platformFamily:  string(rit.PlatformFamily),
		platformVersion: rit.PlatformVersion,
	}
	return targetMap
}
