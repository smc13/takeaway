package services

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

type RedisService struct{}

func (s *RedisService) GetName() string {
	return "Redis"
}

func (s *RedisService) GetDefaultPort() int {
	return 6379
}

func (s *RedisService) GetOrganization() string {
	return ""
}

func (s *RedisService) GetImageName() string {
	return "redis"
}

func (s *RedisService) GetDefaults() map[string]string {
	values := map[string]string{
		"volume": "expose_data",
	}

	// merge base defaults with service defaults
	for key, value := range DefaultOptions() {
		values[key] = value
	}

	return values
}

func (s *RedisService) Prompt() (map[string]string, error) {
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

func (s *RedisService) GetDockerCommandArgs(options map[string]string) []string {
	return []string{
		fmt.Sprintf("--publish=%s:6379", options["port"]),
		fmt.Sprintf("--volume=%s:/data", options["volume"]),
	}
}
