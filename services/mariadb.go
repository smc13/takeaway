package services

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

type MariaDBService struct{}

func (s *MariaDBService) GetName() string {
	return "MariaDB"
}

func (s *MariaDBService) GetDefaultPort() int {
	return 3306
}

func (s *MariaDBService) GetOrganization() string {
	return ""
}

func (s *MariaDBService) GetImageName() string {
	return "mariadb"
}

func (s *MariaDBService) GetDefaults() map[string]string {
	values := map[string]string{
		"volume":        "mariadb_data",
		"root_password": "password",
	}
	// merge base defaults with service defaults
	for key, value := range DefaultOptions() {
		values[key] = value
	}

	return values
}

func (s *MariaDBService) Prompt() (map[string]string, error) {
	defaults := s.GetDefaults()

	prompts := []*survey.Question{
		{
			Name:     "volume",
			Prompt:   &survey.Input{Message: "What is the Docker volume name?", Default: defaults["volume"]},
			Validate: survey.Required,
		},
		{
			Name:     "root_password",
			Prompt:   &survey.Input{Message: "What will the root password be?", Default: defaults["root_password"]},
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

func (s *MariaDBService) GetDockerCommandArgs(options map[string]string) []string {
	return []string{
		fmt.Sprintf("--publish=%s:3306", options["port"]),
		fmt.Sprintf("--volume=%s:/var/lib/mysql", options["volume"]),
		fmt.Sprintf("--env=MYSQL_ROOT_PASSWORD=%s", options["password"]),
	}
}
