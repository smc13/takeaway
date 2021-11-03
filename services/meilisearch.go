package services

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

type MeiliSearchService struct{}

func (s *MeiliSearchService) GetName() string {
	return "MeiliSearch"
}

func (s *MeiliSearchService) GetDefaultPort() int {
	return 7700
}

func (s *MeiliSearchService) GetOrganization() string {
	return "getmeili"
}

func (s *MeiliSearchService) GetImageName() string {
	return "meilisearch"
}

func (s *MeiliSearchService) GetDefaults() map[string]string {
	values := map[string]string{
		"volume": "meili_data",
	}

	// merge base defaults with service defaults
	for key, value := range DefaultOptions() {
		values[key] = value
	}

	return values
}

func (s *MeiliSearchService) Prompt() (map[string]string, error) {
	defaults := s.GetDefaults()

	prompts := []*survey.Question{
		{
			Name:     "volume",
			Prompt:   &survey.Input{Message: "What is the Docker volume name?", Default: defaults["volume"]},
			Validate: survey.Required,
		},
	}

	prompts = append(DefaultPrompts(s.GetDefaultPort()), prompts...)

	var answers = make(map[string]interface{})
	err := survey.Ask(prompts, &answers)
	if err != nil {
		os.Exit(1)
	}

	var mapped = make(map[string]string)
	// convert all answers into strings and add to mapped
	for key, value := range answers {
		mapped[key] = value.(string)
	}

	return mapped, nil
}

func (s *MeiliSearchService) GetDockerCommandArgs(options map[string]string) []string {
	return []string{
		fmt.Sprintf("--publish=%s:7700", options["port"]),
		fmt.Sprintf("--volume=%s:/data.ms", options["volume"]),
	}
}
