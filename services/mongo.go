package services

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

type MongoService struct{}

func (s *MongoService) GetName() string {
	return "Mongo"
}

func (s *MongoService) GetDefaultPort() int {
	return 27017
}

func (s *MongoService) GetOrganization() string {
	return ""
}

func (s *MongoService) GetImageName() string {
	return "mongo"
}

func (s *MongoService) GetDefaults() map[string]string {
	values := map[string]string{
		"volume":        "mongo_data",
		"root_user":     "admin",
		"root_password": "password",
	}
	// merge base defaults with service defaults
	for key, value := range DefaultOptions() {
		values[key] = value
	}

	return values
}

func (s *MongoService) Prompt() (map[string]string, error) {
	defaults := s.GetDefaults()

	prompts := []*survey.Question{
		{
			Name:     "volume",
			Prompt:   &survey.Input{Message: "What is the Docker volume name?", Default: defaults["volume"]},
			Validate: survey.Required,
		},
		{
			Name:     "root_user",
			Prompt:   &survey.Input{Message: "What will the root user be?", Default: defaults["root_user"]},
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

func (s *MongoService) GetDockerCommandArgs(options map[string]string) []string {
	return []string{
		fmt.Sprintf("--publish=%s:27017", options["port"]),
		fmt.Sprintf("--volume=%s:/data/db", options["volume"]),
		fmt.Sprintf("--env=MONGO_INITDB_ROOT_USERNAME=%s", options["root_user"]),
		fmt.Sprintf("--env=MONGO_INITDB_ROOT_PASSWORD=%s", options["root_password"]),
	}
}
