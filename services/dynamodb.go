package services

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

type DynamoDBService struct{}

func (s *DynamoDBService) GetName() string {
	return "DynamoDB"
}

func (s *DynamoDBService) GetDefaultPort() int {
	return 8000
}

func (s *DynamoDBService) GetOrganization() string {
	return "amazon"
}

func (s *DynamoDBService) GetImageName() string {
	return "dynamodb-local"
}

func (s *DynamoDBService) GetDefaults() map[string]string {
	values := map[string]string{
		"volume": "dynamodb_data",
	}
	// merge base defaults with service defaults
	for key, value := range DefaultOptions() {
		values[key] = value
	}

	return values
}

func (s *DynamoDBService) Prompt() (map[string]string, error) {
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

func (s *DynamoDBService) GetDockerCommandArgs(options map[string]string) []string {
	return []string{
		fmt.Sprintf("--publish=%s:8000", options["port"]),
		fmt.Sprintf("--volume=%s:/data.ms", options["volume"]),
		"-u root",
		"jar DynamobDBLocal.jar --sharedDb -dbPath /dynamodb_local_db",
}
}