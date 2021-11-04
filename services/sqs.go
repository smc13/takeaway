package services

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

type SqsService struct{}

func (s *SqsService) GetName() string {
	return "Sqs"
}

func (s *SqsService) GetDefaultPort() int {
	return 9324
}

func (s *SqsService) GetOrganization() string {
	return "roribio16"
}

func (s *SqsService) GetImageName() string {
	return "alpine-sqs"
}

func (s *SqsService) GetDefaults() map[string]string {
	values := map[string]string{
		"volume":          "sqs_data",
		"management_port": "9325",
	}
	// merge base defaults with service defaults
	for key, value := range DefaultOptions() {
		values[key] = value
	}

	return values
}

func (s *SqsService) Prompt() (map[string]string, error) {
	defaults := s.GetDefaults()

	prompts := []*survey.Question{
		{
			Name:     "volume",
			Prompt:   &survey.Input{Message: "What is the Docker volume name?", Default: defaults["volume"]},
			Validate: survey.Required,
		},
		{
			Name:     "management_port",
			Prompt:   &survey.Input{Message: "What will the management port be?", Default: defaults["management_port"]},
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

func (s *SqsService) GetDockerCommandArgs(options map[string]string) []string {
	return []string{
		fmt.Sprintf("--publish=%s:9324", options["port"]),
		fmt.Sprintf("--volume=%s:/data", options["volume"]),
		fmt.Sprintf("--publish=%s:9325", options["management_port"]),
	}
}
