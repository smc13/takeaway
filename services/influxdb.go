package services

import (
	"fmt"
	"os"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
)

type InfluxDBService struct{}

func (s *InfluxDBService) GetName() string {
	return "InfluxDB"
}

func (s *InfluxDBService) GetDefaultPort() int {
	return 8086
}

func (s *InfluxDBService) GetOrganization() string {
	return ""
}

func (s *InfluxDBService) GetImageName() string {
	return "influxdb"
}

func (s *InfluxDBService) GetDefaults() map[string]string {
	values := map[string]string{
		"volume":         "influxdb_data",
		"port":           strconv.Itoa(s.GetDefaultPort()),
		"admin_user":     "admin",
		"admin_password": "password",
	}
	// merge base defaults with service defaults
	for key, value := range DefaultOptions() {
		values[key] = value
	}

	return values
}

func (s *InfluxDBService) Prompt() (map[string]string, error) {
	defaults := s.GetDefaults()

	prompts := []*survey.Question{
		{
			Name:     "volume",
			Prompt:   &survey.Input{Message: "What is the Docker volume name?", Default: defaults["volume"]},
			Validate: survey.Required,
		},
		{
			Name:     "admin_user",
			Prompt:   &survey.Input{Message: "What will the admin user be called?", Default: defaults["admin_user"]},
			Validate: survey.Required,
		},
		{
			Name:     "admin_password",
			Prompt:   &survey.Input{Message: "What will the admin password be?", Default: defaults["admin_password"]},
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

func (s *InfluxDBService) GetDockerCommandArgs(options map[string]string) []string {
	return []string{
		fmt.Sprintf("--publish=%s:8086", options["port"]),
		fmt.Sprintf("--volume=%s:/var/lib/influxdb", options["volume"]),
		"-e",
		fmt.Sprintf("INFLUXDB_ADMIN_USER=%s", options["admin_user"]),
		"-e",
		fmt.Sprintf("INFLUXDB_ADMIN_PASSWORD=%s", options["admin_password"]),
	}
}
