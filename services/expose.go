package services

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

type ExposeService struct{}

func (s *ExposeService) GetName() string {
	return "Expose"
}

func (s *ExposeService) GetDefaultPort() int {
	return 8080
}

func (s *ExposeService) GetOrganization() string {
	return "beyondcodegmbh"
}

func (s *ExposeService) GetImageName() string {
	return "expose-server"
}

func (s *ExposeService) GetDefaults() map[string]string {
	values := map[string]string{
		"domain":   "example.com",
		"volume":   "expose_data",
		"username": "admin",
		"password": "password",
	}

	// merge base defaults with service defaults
	for key, value := range DefaultOptions() {
		values[key] = value
	}

	return values
}

func (s *ExposeService) Prompt() (map[string]string, error) {
	defaults := s.GetDefaults()

	prompts := []*survey.Question{
		{
			Name:     "domain",
			Prompt:   &survey.Input{Message: "What is the domain?", Default: defaults["domain"]},
			Validate: survey.Required,
		},
		{
			Name:     "volume",
			Prompt:   &survey.Input{Message: "What is the Docker volume name?", Default: defaults["volume"]},
			Validate: survey.Required,
		},
		{
			Name:     "username",
			Prompt:   &survey.Input{Message: "What is the username?", Default: defaults["username"]},
			Validate: survey.Required,
		},
		{
			Name:     "password",
			Prompt:   &survey.Input{Message: "What is the password?", Default: defaults["password"]},
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

func (s *ExposeService) GetDockerCommandArgs(options map[string]string) []string {
	return []string{
		fmt.Sprintf("--publish=%s:8080", options["port"]),
		fmt.Sprintf("--volume=%s:/root/.express", options["volume"]),
		fmt.Sprintf("--env=\"username=%s\"", options["username"]),
		fmt.Sprintf("--env=\"password=%s\"", options["password"]),
		fmt.Sprintf("--env=\"domain=%s\"", options["domain"]),
	}
}
