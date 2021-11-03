package services

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

type SingleStoreService struct{}

func (s *SingleStoreService) GetName() string {
	return "SingleStore"
}

func (s *SingleStoreService) GetDefaultPort() int {
	return 3306
}

func (s *SingleStoreService) GetOrganization() string {
	return "singlestore"
}

func (s *SingleStoreService) GetImageName() string {
	return "cluster-in-a-box"
}

func (s *SingleStoreService) GetDefaults() map[string]string {
	values := map[string]string{
		"volume": "singlestore_data",
		"http_port": "8080",
		"license": "",
		"root_password": "password",
	}
	// merge base defaults with service defaults
	for key, value := range DefaultOptions() {
		values[key] = value
	}

	return values
}

func (s *SingleStoreService) Prompt() (map[string]string, error) {
	defaults := s.GetDefaults()

	prompts := []*survey.Question{
		{
			Name:     "volume",
			Prompt:   &survey.Input{Message: "What is the Docker volume name?", Default: defaults["volume"]},
			Validate: survey.Required,
		},
		{
			Name:     "http_port",
			Prompt:   &survey.Input{Message: "Which http host port would you like to use?", Default: defaults["http_port"]},
			Validate: survey.Required,
		},
		{
			Name:     "license",
			Prompt:   &survey.Input{Message: "What is your license key?", Default: defaults["license"]},
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

func (s *SingleStoreService) GetDockerCommandArgs(options map[string]string) []string {
	return []string{
		fmt.Sprintf("--publish=%s:3306", options["port"]),
		fmt.Sprintf("--volume=%s:/var/lib/memsql", options["volume"]),
		fmt.Sprintf("--publish=%s:8080", options["http_port"]),
		fmt.Sprintf("--env=LICENSE_KEY=%s", options["license"]),
		fmt.Sprintf("--env=ROOT_PASSWORD=%s", options["root_password"]),
	}
}