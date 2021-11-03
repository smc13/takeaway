package services

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

type CouchDbService struct{}

func (s *CouchDbService) GetName() string {
	return "CouchDB"
}

func (s *CouchDbService) GetDefaultPort() int {
	return 5984
}

func (s *CouchDbService) GetOrganization() string {
	return ""
}

func (s *CouchDbService) GetImageName() string {
	return "couchdb"
}

func (s *CouchDbService) GetDefaults() map[string]string {
	values := map[string]string{
		"volume": "couchdb_data",
	}

	// merge base defaults with service defaults
	for key, value := range DefaultOptions() {
		values[key] = value
	}

	return values
}

func (s *CouchDbService) Prompt() (map[string]string, error) {
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

func (s *CouchDbService) GetDockerCommandArgs(options map[string]string) []string {
	return []string{
		fmt.Sprintf("--publish=%s:5984", options["port"]),
		fmt.Sprintf("--volume=%s:/opt/couchdb/data", options["volume"]),
	}
}
