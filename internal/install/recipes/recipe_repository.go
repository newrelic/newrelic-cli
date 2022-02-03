package recipes

import (
	"math"
	"regexp"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

var (
	kernelArch      = "KernelArch"
	kernelVersion   = "KernelVersion"
	oS              = "OS"
	platform        = "Platform"
	platformFamily  = "PlatformFamily"
	platformVersion = "PlatformVersion"
)

type RecipeRepository struct {
	RecipeLoaderFunc  func() ([]types.OpenInstallationRecipe, error)
	loadedRecipes     []types.OpenInstallationRecipe
	filteredRecipes   []types.OpenInstallationRecipe
	discoveryManifest *types.DiscoveryManifest
}

type recipeMatch struct {
	matchCount int
	recipe     types.OpenInstallationRecipe
}

// NewRecipeRepository returns a new instance of types.RecipeRepository.
func NewRecipeRepository(loaderFunc func() ([]types.OpenInstallationRecipe, error), manifest *types.DiscoveryManifest) *RecipeRepository {
	rr := RecipeRepository{
		RecipeLoaderFunc:  loaderFunc,
		loadedRecipes:     nil,
		filteredRecipes:   nil,
		discoveryManifest: manifest,
	}

	return &rr
}

func (rf *RecipeRepository) FindRecipeByName(name string) *types.OpenInstallationRecipe {
	recipes, err := rf.FindAll()
	if err != nil {
		log.Error(err)
		return nil
	}

	for _, r := range recipes {
		if strings.EqualFold(r.Name, name) {
			log.Debugf("Found recipe name %v", r.Name)
			return &r
		}
	}
	return nil
}

func (rf *RecipeRepository) FindAll() ([]types.OpenInstallationRecipe, error) {
	if rf.filteredRecipes != nil {
		return rf.filteredRecipes, nil
	}

	if rf.loadedRecipes == nil {
		recipes, err := rf.RecipeLoaderFunc()
		if err != nil {
			return nil, err
		}
		log.Debugf("Loaded %d recipes", len(recipes))
		rf.loadedRecipes = recipes
	}

	rf.filteredRecipes = filterRecipes(rf.loadedRecipes, rf.discoveryManifest)

	return rf.filteredRecipes, nil
}

func filterRecipes(loadedRecipes []types.OpenInstallationRecipe, m *types.DiscoveryManifest) []types.OpenInstallationRecipe {
	results := []types.OpenInstallationRecipe{}
	matchRecipes := make(map[string][]recipeMatch)
	hostMap := getHostMap(m)

	log.Debugf("Find all available out of %d recipes for host %+v", len(loadedRecipes), hostMap)

	for _, recipe := range loadedRecipes {
		matchTargetCount := []int{}

		for _, rit := range recipe.InstallTargets {
			matchCount := 0
			for k, v := range getRecipeTargetMap(rit) {
				if v == "" {
					continue
				}
				isValueMatching := matchRecipeCriteria(hostMap, k, v)
				if isValueMatching {
					log.Tracef("matching recipe %s field name %s and value %s using hostMap %+v", recipe.Name, k, v, hostMap)
					matchCount++
				} else {
					log.Tracef("recipe %s defines %s=%s but input did not provide a match using hostMap %+v", recipe.Name, k, v, hostMap)
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
			log.Tracef("Recipe InstallTargetsCount %d and maxMatchCount %d", len(recipe.InstallTargets), maxMatchTargetCount)

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

	if len(matchRecipes) > 0 {
		keys := []string{}
		unorderedResults := map[string]types.OpenInstallationRecipe{}
		for _, matches := range matchRecipes {
			if len(matches) > 0 {
				match := findMaxMatch(matches)
				singleRecipe := match.recipe
				log.Tracef("Add result for recipe name %s with targets %+v", match.recipe, match.recipe.InstallTargets)
				key := singleRecipe.GetOrderKey()
				keys = append(keys, key)
				unorderedResults[key] = singleRecipe
			}
		}
		sort.Strings(keys)
		for _, k := range keys {
			recipe := unorderedResults[k]
			results = append(results, recipe)
		}
	}
	return results
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
		if len(rvalue) > 0 && rvalue[0] == '(' {
			if regex, error := regexp.Compile(rvalue); error == nil {
				return regex.MatchString(val)
			}
		}
		return strings.EqualFold(val, rvalue)
	}

	return false
}

func getHostMap(m *types.DiscoveryManifest) map[string]string {
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
