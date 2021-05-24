package types

func (r RecipeVars) ToSlice() []string {
	var s []string
	for k, v := range r {
		s = append(s, []string{k, v}...)
	}

	return s
}
