package install

import "github.com/manifoldco/promptui"

type promptUIPrompter struct{}

func (p *promptUIPrompter) promptYesNo(msg string) (bool, error) {
	prompt := promptui.Select{
		Label: msg,
		Items: []string{"Yes", "No"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return false, err
	}

	return result == "Yes", nil
}
