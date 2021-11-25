package services

import (
	"fmt"
	"os"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
)

type BeanstalkdService struct{}

func (s *BeanstalkdService) GetName() string {
	return "Beanstalkd"
}

func (s *BeanstalkdService) GetOrganization() string {
	return "schickling"
}

func (s *BeanstalkdService) GetImageName() string {
	return "beanstalkd"
}

func (s *BeanstalkdService) GetDefaultPort() int {
	return 11300
}

func (s *BeanstalkdService) GetDefaults() map[string]string {
	values := map[string]string{
		"port": strconv.Itoa(s.GetDefaultPort()),
	}

	// merge base defaults with service defaults
	for key, value := range DefaultOptions() {
		values[key] = value
	}

	return values
}

func (s *BeanstalkdService) Prompt() (map[string]string, error) {
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

func (s *BeanstalkdService) GetDockerCommandArgs(options map[string]string) []string {
	return []string{
		fmt.Sprintf("--publish=%s:6379", options["port"]),
	}
}
