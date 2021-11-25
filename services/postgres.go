package services

import (
	"fmt"
	"os"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
)

type PostgresService struct{}

func (s *PostgresService) GetName() string {
	return "PostgresSQL"
}

func (s *PostgresService) GetDefaultPort() int {
	return 5432
}

func (s *PostgresService) GetOrganization() string {
	return ""
}

func (s *PostgresService) GetImageName() string {
	return "postgres"
}

func (s *PostgresService) GetDefaults() map[string]string {
	values := map[string]string{
		"volume":   "postgres_data",
		"port":     strconv.Itoa(s.GetDefaultPort()),
		"password": "password",
	}

	// merge base defaults with service defaults
	for key, value := range DefaultOptions() {
		values[key] = value
	}

	return values
}

func (s *PostgresService) Prompt() (map[string]string, error) {
	defaults := s.GetDefaults()

	prompts := []*survey.Question{
		{
			Name:     "volume",
			Prompt:   &survey.Input{Message: "What is the Docker volume name?", Default: defaults["volume"]},
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

func (s *PostgresService) GetDockerCommandArgs(options map[string]string) []string {
	return []string{
		fmt.Sprintf("--publish=%s:5432", options["port"]),
		fmt.Sprintf("--volume=%s:/var/lib/postgresql/data", options["volume"]),
		"-e",
		fmt.Sprintf("\"POSTGRES_PASSWORD=%s\"", options["password"]),
	}
}
