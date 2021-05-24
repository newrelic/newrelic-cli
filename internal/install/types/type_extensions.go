package types

import "fmt"

func (r RecipeVars) ToSlice() []string {
	var s []string
	for k, v := range r {
		s = append(s, fmt.Sprintf("%s=%s", k, v))
	}

	return s
}
