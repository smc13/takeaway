package services

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

type MailDevService struct{}

func (s *MailDevService) GetName() string {
	return "MailDev"
}

func (s *MailDevService) GetDefaultPort() int {
	return 1025
}

func (s *MailDevService) GetOrganization() string {
	return "maildev"
}

func (s *MailDevService) GetImageName() string {
	return "maildev"
}

func (s *MailDevService) GetDefaults() map[string]string {
	values := map[string]string{
		"web_port": "8025",
	}
	// merge base defaults with service defaults
	for key, value := range DefaultOptions() {
		values[key] = value
	}

	return values
}

func (s *MailDevService) Prompt() (map[string]string, error) {
	defaults := s.GetDefaults()

	prompts := []*survey.Question{
		DefaultPrompts(s.GetDefaultPort())[0],
		{
			Name:   "tag",
			Prompt: &survey.Input{Message: "Which tag (version) would you like to use?", Default: "1.1.0"},
		},
		{
			Name: "web_port",
			Prompt: &survey.Input{
				Message: "What will the web port be?",
				Default: defaults["web_port"],
			},
			Validate: survey.Required,
		},
	}

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

func (s *MailDevService) GetDockerCommandArgs(options map[string]string) []string {
	return []string{
		fmt.Sprintf("--publish=%s:25", options["port"]),
		fmt.Sprintf("--publish=%s:80", options["web_port"]),
	}
}
