package services

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

type MemcachedService struct{}

func (s *MemcachedService) GetName() string {
	return "Memcached"
}

func (s *MemcachedService) GetDefaultPort() int {
	return 11211
}

func (s *MemcachedService) GetOrganization() string {
	return ""
}

func (s *MemcachedService) GetImageName() string {
	return "memcached"
}

func (s *MemcachedService) GetDefaults() map[string]string {
	values := map[string]string{}
	// merge base defaults with service defaults
	for key, value := range DefaultOptions() {
		values[key] = value
	}

	return values
}

func (s *MemcachedService) Prompt() (map[string]string, error) {
	prompts := DefaultPrompts(s.GetDefaultPort())

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

func (s *MemcachedService) GetDockerCommandArgs(options map[string]string) []string {
	return []string{
		fmt.Sprintf("--publish=%s:11211", options["port"]),
	}
}
