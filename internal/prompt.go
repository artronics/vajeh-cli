package internal

import (
	"github.com/manifoldco/promptui"
)

type PromptData struct {
	Label        string
	DefaultValue string
	Key          string
}

func GetPromptResult(prompts []PromptData) (map[string]string, error) {
	m := make(map[string]string, len(prompts))

	for _, p := range prompts {
		res, err := runPrompt(p)
		if err != nil {
			return m, err
		}

		m[p.Key] = res
	}

	return m, nil
}

func runPrompt(p PromptData) (string, error) {
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}
	prompt := promptui.Prompt{
		Label:     p.Label,
		Default:   p.DefaultValue,
		Templates: templates,
	}

	return prompt.Run()
}
