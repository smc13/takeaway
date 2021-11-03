package services

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/smcassar/takeaway/docker"
	"golang.org/x/mod/semver"
)

var services = []Service{
	&BeanstalkdService{},
	&ClickHouseService{},
	&CouchDbService{},
	&DynamoDBService{},
	&ElasticSearchService{},
	&EventStoreDbService{},
	&ExposeService{},
	&InfluxDBService{},
	&MailDevService{},
	&MailHogService{},
	&MariaDBService{},
	&MeiliSearchService{},
	&MemcachedService{},
	&MinioService{},
	&MongoService{},
	&MsSqlService{},
	&MySqlService{},
	&Neo4jService{},
	&PostGisService{},
	&PostgresService{},
	&RabbitMqService{},
	&RedisService{},
	&SftpService{},
	&SingleStoreService{},
	&SqsService{},
}

// Get all available services
func GetServices() []Service {
	return services
}

// Get the names of all enabled services
func GetServiceNames() []string {
	names := make([]string, len(services))
	for i, s := range GetServices() {
		names[i] = s.GetName()
	}

	return names
}

// Get a service by name
func GetService(name string) Service {
	for _, s := range services {
		if strings.EqualFold(s.GetName(), name) {
			return s
		}
	}

	return nil
}

func EnableService(service Service, useDefaults bool) {
	var options map[string]string

	if useDefaults {
		options = service.GetDefaults()
		color.Yellow("Using default options for %s", service.GetName())
	} else {
		color.Yellow("\nConfiguring %s", service.GetName())
		answers, err := service.Prompt()
		if err != nil {
			color.Red("Error: %s", err)
			return
		}

		options = answers
	}

	err := docker.EnsureNetworkCreated()
	if err != nil {
		color.Red("Error while creating network: %s", err)
		return
	}

	tag, err := docker.ResolveTag(service.GetOrganization(), service.GetImageName(), options["tag"])
	if err != nil {
		color.Red("Error while resolving tag: %s", err)
		return
	}

	// ensure the image is downloaded
	imageErr := ensureImageIsDownloaded(service, tag)
	if imageErr != nil {
		color.Red("Error while downloading image: %s", imageErr)
		return
	}

	alias := getAlias(service.GetName(), tag)
	options["tag"] = tag
	dockerTemplate := service.GetDockerCommandArgs(options)

	// check that the port is not already in use
	if PortIsInUse(options["port"]) {
		color.Red("Port %s is already in use", options["port"])
		return
	}

	color.Green("Enabling %s...", service.GetName())

	args := append([]string{fmt.Sprintf("--name=%s", getContainerName(service.GetName(), tag, options["port"]))}, docker.GetNetworkSettings(alias, service.GetImageName())...)
	args = append(args, dockerTemplate...)
	args = append(args, buildImageName(service, tag))

	err = docker.CreateContainer(args)
	if err != nil {
		color.Red("Error while creating container: %s", err)
		return
	}

	color.Green("✔ %s enabled", service.GetName())
}

// Generate a container name based on the service name, tag & port
func getContainerName(serviceName string, tag string, port string) string {
	return fmt.Sprintf("TO--%s--%s--%s", strings.ToLower(serviceName), tag, port)
}

func buildImageName(service Service, tag string) string {
	organization := service.GetOrganization()
	if organization == "" {
		organization = "library"
	}

	return fmt.Sprintf("%s/%s:%s", organization, service.GetImageName(), tag)
}

// Generate an alias based on the service name & tag
func getAlias(serviceName string, tag string) string {
	shortName := strings.ToLower(serviceName)
	// check if tag represents a semver version
	if !semver.IsValid(tag) {
		return fmt.Sprintf("%s-%s", shortName, tag)
	}

	return fmt.Sprintf("%s%s", shortName, semver.MajorMinor(tag))
}

// Ensure that the services image is downloaded
func ensureImageIsDownloaded(service Service, tag string) error {
	if !docker.ImageExists(service.GetOrganization(), service.GetImageName(), tag) {
		color.Yellow("Downloading %s/%s:%s...", service.GetOrganization(), service.GetImageName(), tag)

		err := docker.PullImage(service.GetOrganization(), service.GetImageName(), tag)

		if err != nil {
			return err
		}

		color.Green("✔ %s/%s:%s downloaded", service.GetOrganization(), service.GetImageName(), tag)
	}

	return nil
}
