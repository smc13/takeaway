package services

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

type MinioService struct{}

func (s *MinioService) GetName() string {
	return "Minio"
}

func (s *MinioService) GetDefaultPort() int {
	return 9000
}

func (s *MinioService) GetOrganization() string {
	return "minio"
}

func (s *MinioService) GetImageName() string {
	return "minio"
}

func (s *MinioService) GetDefaults() map[string]string {
	values := map[string]string{
		"volume":        "minio_data",
		"console":       "9001",
		"root_user":     "minioadmin",
		"root_password": "minioadmin",
	}
	// merge base defaults with service defaults
	for key, value := range DefaultOptions() {
		values[key] = value
	}

	return values
}

func (s *MinioService) Prompt() (map[string]string, error) {
	defaults := s.GetDefaults()

	prompts := []*survey.Question{
		{
			Name:     "volume",
			Prompt:   &survey.Input{Message: "What is the Docker volume name?", Default: defaults["volume"]},
			Validate: survey.Required,
		},
		{
			Name:     "console",
			Prompt:   &survey.Input{Message: "Which host port would you like to be used by Minio Console?", Default: defaults["console"]},
			Validate: survey.Required,
		},
		{
			Name:     "root_user",
			Prompt:   &survey.Input{Message: "What will the root user name for Minio be?", Default: defaults["root_user"]},
			Validate: survey.Required,
		},
		{
			Name:     "root_password",
			Prompt:   &survey.Input{Message: "What will the root password for Minio be?", Default: defaults["root_password"]},
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

func (s *MinioService) GetDockerCommandArgs(options map[string]string) []string {
	return []string{
		fmt.Sprintf("--publish=%s:9000", options["port"]),
		fmt.Sprintf("--volume=%s:/data", options["volume"]),
		fmt.Sprintf("--publish=%s:9001", options["console"]),
		fmt.Sprintf("--env=MINIO_ROOT_USER=%s", options["root_user"]),
		fmt.Sprintf("--env=MINIO_ROOT_PASSWORD=%s", options["root_password"]),
	}
}
