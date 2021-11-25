package services

import (
	"fmt"
	"os"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
)

type MailHogService struct{}

func (s *MailHogService) GetName() string {
	return "MailHog"
}

func (s *MailHogService) GetDefaultPort() int {
	return 1025
}

func (s *MailHogService) GetOrganization() string {
	return "mailhog"
}

func (s *MailHogService) GetImageName() string {
	return "mailhog"
}

func (s *MailHogService) GetDefaults() map[string]string {
	values := map[string]string{
		"port":     strconv.Itoa(s.GetDefaultPort()),
		"web_port": "8025",
	}

	// merge base defaults with service defaults
	for key, value := range DefaultOptions() {
		values[key] = value
	}

	return values
}

func (s *MailHogService) Prompt() (map[string]string, error) {
	defaults := s.GetDefaults()

	prompts := []*survey.Question{
		{
			Name:     "web_port",
			Prompt:   &survey.Input{Message: "What will the web port be?", Default: defaults["web_port"]},
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

func (s *MailHogService) GetDockerCommandArgs(options map[string]string) []string {
	return []string{
		fmt.Sprintf("--publish=%s:1025", options["port"]),
		fmt.Sprintf("--publish=%s:8025", options["port"]),
	}
}
