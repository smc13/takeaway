package services

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

type ClickHouseService struct{}

func (s *ClickHouseService) GetName() string {
	return "ClickHouse"
}

func (s *ClickHouseService) GetDefaultPort() int {
	return 9000
}

func (s *ClickHouseService) GetOrganization() string {
	return "yandex"
}

func (s *ClickHouseService) GetImageName() string {
	return "clickhouse-server"
}

func (s *ClickHouseService) GetDefaults() map[string]string {
	values := map[string]string{
		"http_port": "8123",
		"volume":    "clickhouse_data",
	}

	// merge base defaults with service defaults
	for key, value := range DefaultOptions() {
		values[key] = value
	}

	return values
}

func (s *ClickHouseService) Prompt() (map[string]string, error) {
	defaults := s.GetDefaults()

	prompts := []*survey.Question{
		{
			Name:     "http_port",
			Prompt:   &survey.Input{Message: "Which http host port would you like to use?", Default: defaults["http_port"]},
			Validate: survey.Required,
		},
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

func (s *ClickHouseService) GetDockerCommandArgs(options map[string]string) []string {
	return []string{
		fmt.Sprintf("--publish=%s:9000", options["port"]),
		fmt.Sprintf("--publish=%s:8123", options["http_port"]),
		fmt.Sprintf("--volume=%s:/var/lib/clickhouse", options["volume"]),
		"--ulimit=nofile=262144:26144",
	}
}
