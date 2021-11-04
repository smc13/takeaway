package services

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

type MsSqlService struct{}

func (s *MsSqlService) GetName() string {
	return "MS SQL Server"
}

func (s *MsSqlService) GetDefaultPort() int {
	return 1433
}

func (s *MsSqlService) GetOrganization() string {
	return "mcr.microsoft.com"
}

func (s *MsSqlService) GetImageName() string {
	return "mssql/server"
}

func (s *MsSqlService) GetDefaults() map[string]string {
	values := map[string]string{
		"sa_password": "useA$strongPas1337",
	}
	// merge base defaults with service defaults
	for key, value := range DefaultOptions() {
		values[key] = value
	}

	return values
}

func (s *MsSqlService) Prompt() (map[string]string, error) {
	defaults := s.GetDefaults()

	prompts := []*survey.Question{
		{
			Name:     "sa_password",
			Prompt:   &survey.Input{Message: "What wil the password for the `sa` user be?", Default: defaults["sa_password"]},
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

func (s *MsSqlService) GetDockerCommandArgs(options map[string]string) []string {
	return []string{
		fmt.Sprintf("--publish=%s:1433", options["port"]),
		"--env=ACCEPT_EULA=Y",
		fmt.Sprintf("--env=SA_PASSWORD=%s", options["sa_password"]),
	}
}
