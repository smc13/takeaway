package services

import (
	"fmt"
	"os"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
)

type ElasticSearchService struct{}

func (s *ElasticSearchService) GetName() string {
	return "Elasticsearch"
}

func (s *ElasticSearchService) GetDefaultPort() int {
	return 9200
}

func (s *ElasticSearchService) GetOrganization() string {
	return "docker.elastic.co"
}

func (s *ElasticSearchService) GetImageName() string {
	return "elasticsearch/elasticsearch"
}

func (s *ElasticSearchService) GetDefaults() map[string]string {
	values := map[string]string{
		"volume": "elastic_data",
		"port":   strconv.Itoa(s.GetDefaultPort()),
	}

	// merge base defaults with service defaults
	for key, value := range DefaultOptions() {
		values[key] = value
	}

	return values
}

func (s *ElasticSearchService) Prompt() (map[string]string, error) {
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

func (s *ElasticSearchService) GetDockerCommandArgs(options map[string]string) []string {
	return []string{
		fmt.Sprintf("--publish=%s:9200", options["port"]),
		fmt.Sprintf("--volume=%s:/usr/share/elasticsearch/data", options["volume"]),
		"-e",
		"discovery.type=single-node",
	}
}
