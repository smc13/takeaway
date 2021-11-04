package services

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

type RabbitMqService struct{}

func (s *RabbitMqService) GetName() string {
	return "RabbitMQ"
}

func (s *RabbitMqService) GetDefaultPort() int {
	return 5672
}

func (s *RabbitMqService) GetOrganization() string {
	return "rabbitmq"
}

func (s *RabbitMqService) GetImageName() string {
	return "rabbitmq"
}

func (s *RabbitMqService) GetDefaults() map[string]string {
	values := map[string]string{
		"hostname":  "takeout",
		"mgmt_port": "15672",
	}
	// merge base defaults with service defaults
	for key, value := range DefaultOptions() {
		values[key] = value
	}

	return values
}

func (s *RabbitMqService) Prompt() (map[string]string, error) {
	defaults := s.GetDefaults()

	prompts := []*survey.Question{
		DefaultPrompts(s.GetDefaultPort())[0],
		{
			Name:   "tag",
			Prompt: &survey.Input{Message: "Which tag (version) would you like to use?", Default: "management"},
		},
		{
			Name:     "hostname",
			Prompt:   &survey.Input{Message: "Which hostname would you like to user?", Default: defaults["hostname"]},
			Validate: survey.Required,
		},
		{
			Name:     "mgmt_port",
			Prompt:   &survey.Input{Message: "Which management port would you like to use?", Default: defaults["mgmt_port"]},
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

func (s *RabbitMqService) GetDockerCommandArgs(options map[string]string) []string {
	return []string{
		fmt.Sprintf("--publish=%s:5672", options["port"]),
		fmt.Sprintf("--publish=%s:15672", options["mgmt_port"]),
		fmt.Sprintf("--hostname=%s", options["hostname"]),
	}
}
