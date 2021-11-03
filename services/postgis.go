package services

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

type PostGisService struct{}

func (s *PostGisService) GetName() string {
	return "PostGis"
}

func (s *PostGisService) GetDefaultPort() int {
	return 5432
}

func (s *PostGisService) GetOrganization() string {
	return "postgis"
}

func (s *PostGisService) GetImageName() string {
	return "postgis"
}

func (s *PostGisService) GetDefaults() map[string]string {
	values := map[string]string{
		"volume": "postgis_data",
		"root_password": "password",
	}
	// merge base defaults with service defaults
	for key, value := range DefaultOptions() {
		values[key] = value
	}

	return values
}

func (s *PostGisService) Prompt() (map[string]string, error) {
	defaults := s.GetDefaults()

	prompts := []*survey.Question{
		{
			Name:     "volume",
			Prompt:   &survey.Input{Message: "What is the Docker volume name?", Default: defaults["volume"]},
			Validate: survey.Required,
		},
		{
			Name:     "root_password",
			Prompt:   &survey.Input{Message: "What will the password for the `postgres` user be?", Default: defaults["root_password"]},
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

func (s *PostGisService) GetDockerCommandArgs(options map[string]string) []string {
	return []string{
		fmt.Sprintf("--publish=%s:5432", options["port"]),
		fmt.Sprintf("--volume=%s:/var/lib/postgis/data", options["volume"]),
		fmt.Sprintf("--env=POSTGRES_PASSWORD=%s", options["root_password"]),
	}
}