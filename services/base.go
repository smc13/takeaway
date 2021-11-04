package services

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
)

type Service interface {
	GetName() string
	GetDefaultPort() int
	GetOrganization() string
	GetImageName() string
	GetDefaults() map[string]string
	Prompt() (map[string]string, error)
	GetDockerCommandArgs(map[string]string) []string
}

func DefaultOptions() map[string]string {
	return map[string]string{
		"tag": "latest",
	}
}

func DefaultPrompts(port int) []*survey.Question {
	return []*survey.Question{
		{
			Name:   "port",
			Prompt: &survey.Input{Message: "Which port would you like to use?", Default: strconv.Itoa(port)},
			Validate: func(val interface{}) error {
				port, err := strconv.ParseUint(val.(string), 10, 16)
				if err != nil {
					return errors.New("invalid port")
				}

				if PortIsInUse(uint16(port)) {
					return errors.New("port is already in use")
				}

				return nil
			},
		},
		{
			Name:   "tag",
			Prompt: &survey.Input{Message: "Which tag (version) would you like to use?", Default: "latest"},
		},
	}
}

func PortIsInUse(port uint16) bool {
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))

	if err != nil {
		return false
	}

	conn.Close()
	return true
}
