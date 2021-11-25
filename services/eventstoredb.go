package services

import (
	"fmt"
	"os"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
)

type EventStoreDbService struct{}

func (s *EventStoreDbService) GetName() string {
	return "EventStoreDB"
}

func (s *EventStoreDbService) GetDefaultPort() int {
	return 1113
}

func (s *EventStoreDbService) GetOrganization() string {
	return "eventstore"
}

func (s *EventStoreDbService) GetImageName() string {
	return "eventstore"
}

func (s *EventStoreDbService) GetDefaults() map[string]string {
	values := map[string]string{
		"volume":   "eventstoredb_data",
		"port":     strconv.Itoa(s.GetDefaultPort()),
		"web_port": "2113",
	}
	// merge base defaults with service defaults
	for key, value := range DefaultOptions() {
		values[key] = value
	}

	return values
}

func (s *EventStoreDbService) Prompt() (map[string]string, error) {
	defaults := s.GetDefaults()

	prompts := []*survey.Question{
		DefaultPrompts(s.GetDefaultPort())[0],
		{
			Name:   "tag",
			Prompt: &survey.Input{Message: "Which tag (version) would you like to use?", Default: "5.0.8-xenial"},
		},
		{
			Name:     "volume",
			Prompt:   &survey.Input{Message: "What is the Docker volume name?", Default: defaults["volume"]},
			Validate: survey.Required,
		},
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

func (s *EventStoreDbService) GetDockerCommandArgs(options map[string]string) []string {
	return []string{
		fmt.Sprintf("--publish=%s:8080", options["port"]),
		fmt.Sprintf("--volume=%s:/", options["volume"]),
		fmt.Sprintf("--publish=%s:2113", options["web_port"]),
	}
}
