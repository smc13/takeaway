package services

import (
	"fmt"
	"os"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
)

type Neo4jService struct{}

func (s *Neo4jService) GetName() string {
	return "Neo4j"
}

func (s *Neo4jService) GetDefaultPort() int {
	return 7474
}

func (s *Neo4jService) GetOrganization() string {
	return ""
}

func (s *Neo4jService) GetImageName() string {
	return "neo4j"
}

func (s *Neo4jService) GetDefaults() map[string]string {
	values := map[string]string{
		"volume":           "neo4j_data",
		"port":             strconv.Itoa(s.GetDefaultPort()),
		"bolt_access_port": "7687",
	}
	// merge base defaults with service defaults
	for key, value := range DefaultOptions() {
		values[key] = value
	}

	return values
}

func (s *Neo4jService) Prompt() (map[string]string, error) {
	defaults := s.GetDefaults()

	prompts := []*survey.Question{
		{
			Name:     "volume",
			Prompt:   &survey.Input{Message: "What is the Docker volume name?", Default: defaults["volume"]},
			Validate: survey.Required,
		},
		{
			Name:     "bolt_access_port",
			Prompt:   &survey.Input{Message: "What will the Bolt access port be?", Default: defaults["bolt_access_port"]},
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

func (s *Neo4jService) GetDockerCommandArgs(options map[string]string) []string {
	return []string{
		fmt.Sprintf("--publish=%s:7474", options["port"]),
		fmt.Sprintf("--volume=%s:/data", options["volume"]),
		fmt.Sprintf("--publish=%s:7687", options["bolt_access_port"]),
		"-e",
		"NEO4J_AUTH=none",
		"-e",
		"NEO4J_ACCEPT_LICENSE_AGREEMENT=yes",
	}
}
