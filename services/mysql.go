package services

import (
	"fmt"
	"os"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
)

type MySqlService struct{}

func (s *MySqlService) GetName() string {
	return "Mysql"
}

func (s *MySqlService) GetDefaultPort() int {
	return 3306
}

func (s *MySqlService) GetOrganization() string {
	return ""
}

func (s *MySqlService) GetImageName() string {
	return "mysql"
}

func (s *MySqlService) GetDefaults() map[string]string {
	values := map[string]string{
		"volume":        "mysql_data",
		"root_password": "",
		"port":          strconv.Itoa(s.GetDefaultPort()),
	}
	// merge base defaults with service defaults
	for key, value := range DefaultOptions() {
		values[key] = value
	}

	return values
}

func (s *MySqlService) Prompt() (map[string]string, error) {
	defaults := s.GetDefaults()

	prompts := []*survey.Question{
		{
			Name:     "volume",
			Prompt:   &survey.Input{Message: "What is the Docker volume name?", Default: defaults["volume"]},
			Validate: survey.Required,
		},
		{
			Name:   "root_password",
			Prompt: &survey.Input{Message: "What will the root password be? (null by default)", Default: defaults["root_password"]},
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

func (s *MySqlService) GetDockerCommandArgs(options map[string]string) []string {
	allowEmptyPassword := "no"
	if options["root_password"] == "" {
		allowEmptyPassword = "yes"
	}

	return []string{
		fmt.Sprintf("--publish=%s:3306", options["port"]),
		fmt.Sprintf("--volume=%s:/var/lib/mysql", options["volume"]),
		"-e",
		fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", options["root_password"]),
		"-e",
		fmt.Sprintf("MYSQL_ALLOW_EMPTY_PASSWORD=%s", allowEmptyPassword),
	}
}
