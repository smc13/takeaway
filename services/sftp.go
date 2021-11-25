package services

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

type SftpService struct{}

func (s *SftpService) GetName() string {
	return "Sftp"
}

func (s *SftpService) GetDefaultPort() int {
	return 22
}

func (s *SftpService) GetOrganization() string {
	return "atmoz"
}

func (s *SftpService) GetImageName() string {
	return "sftp"
}

func (s *SftpService) GetDefaults() map[string]string {
	values := map[string]string{
		"port":             strconv.Itoa(s.GetDefaultPort()),
		"user_name":        "foo",
		"password":         "pass",
		"upload_directory": "upload",
		"mapped_directory": "",
	}
	// merge base defaults with service defaults
	for key, value := range DefaultOptions() {
		values[key] = value
	}

	return values
}

func (s *SftpService) Prompt() (map[string]string, error) {
	defaults := s.GetDefaults()

	prompts := []*survey.Question{
		DefaultPrompts(s.GetDefaultPort())[0],
		{
			Name:   "tag",
			Prompt: &survey.Input{Message: "Which tag (version) would you like to use?", Default: "alpine"},
		},
		{
			Name:     "user_name",
			Prompt:   &survey.Input{Message: "What will the default username be?", Default: defaults["user_name"]},
			Validate: survey.Required,
		},
		{
			Name:     "password",
			Prompt:   &survey.Input{Message: "What will the default password be?", Default: defaults["password"]},
			Validate: survey.Required,
		},
		{
			Name:     "upload_directory",
			Prompt:   &survey.Input{Message: "Where will files be uploaded?", Default: defaults["upload_directory"]},
			Validate: survey.Required,
		},
		{
			Name:   "mapped_directory",
			Prompt: &survey.Input{Message: "Which local directory should be mapped inside? (nothing if null)", Default: defaults["mapped_directory"]},
		},
	}

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

func (s *SftpService) GetDockerCommandArgs(options map[string]string) []string {
	args := []string{
		fmt.Sprintf("--publish=%s:22", options["port"]),
	}

	if options["mapped_directory"] != "" {
		localMapping := fmt.Sprintf("%s:/home/%s/%s", strings.Trim(options["mapped_directory"], " "), options["user_name"], options["upload_directory"])
		args = append(args, fmt.Sprintf("--volume=%s:", localMapping))
	} else {
		userConfig := fmt.Sprintf("%s:%s:::%s", options["user_name"], options["password"], options["upload_directory"])
		args = append(args, fmt.Sprintf("--user=%s", userConfig))
	}

	return args
}
