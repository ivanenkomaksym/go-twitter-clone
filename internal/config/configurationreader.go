package config

import (
	_ "embed"
	"encoding/json"
	"log"
	"os"
	"strconv"

	"golang.org/x/oauth2/google"
)

//go:embed appsettings.json
var appsettingsContent []byte

func ReadConfiguration() Configuration {
	configuration := Configuration{}
	json.Unmarshal(appsettingsContent, &configuration)

	// Set default google endpoints
	if authenticationEnableEnvVar := os.Getenv("Authentication:Enable"); authenticationEnableEnvVar != "" {
		log.Println("Overriding Authentication:Enable from environment variable: ", authenticationEnableEnvVar)
		configuration.Authentication.Enable, _ = strconv.ParseBool(authenticationEnableEnvVar)
	}

	if configuration.Authentication.Enable {
		configuration.Authentication.OAuth2.Endpoint = google.Endpoint
	}

	if modeEnvVar := os.Getenv("Mode"); modeEnvVar != "" {
		log.Println("Overriding Mode from environment variable: ", modeEnvVar)
		configuration.Mode, _ = ParseMode(modeEnvVar)
	}

	if projectIdEnvVar := os.Getenv("ProjectId"); projectIdEnvVar != "" {
		log.Println("Overriding ProjectId from environment variable: ", projectIdEnvVar)
		configuration.ProjectId = projectIdEnvVar
	}

	if clientIdEnvVar := os.Getenv("Authentication:OAuth2:ClientID"); clientIdEnvVar != "" {
		log.Println("Overriding Authentication:OAuth2:ClientID from environment variable: ", clientIdEnvVar)
		configuration.Authentication.OAuth2.ClientID = clientIdEnvVar
	}

	if clientSecretEnvVar := os.Getenv("Authentication:OAuth2:ClientSecret"); clientSecretEnvVar != "" {
		log.Println("Overriding Authentication:OAuth2:ClientSecret from environment variable: ", clientSecretEnvVar)
		configuration.Authentication.OAuth2.ClientSecret = clientSecretEnvVar
	}

	if tweetsStorageConnectionStringEnvVar := os.Getenv("TweetsStorage:ConnectionString"); tweetsStorageConnectionStringEnvVar != "" {
		log.Println("Overriding TweetsStorage:ConnectionString from environment variable: ", tweetsStorageConnectionStringEnvVar)
		configuration.TweetsStorage.ConnectionString = tweetsStorageConnectionStringEnvVar
	}

	if feedsStorageConnectionStringEnvVar := os.Getenv("FeedsStorage:ConnectionString"); feedsStorageConnectionStringEnvVar != "" {
		log.Println("Overriding FeedsStorage:ConnectionString from environment variable: ", feedsStorageConnectionStringEnvVar)
		configuration.FeedsStorage.ConnectionString = feedsStorageConnectionStringEnvVar
	}

	if apiServerApplicationUrlStringEnvVar := os.Getenv("ApiServer:ApplicationUrl"); apiServerApplicationUrlStringEnvVar != "" {
		log.Println("Overriding ApiServer:ApplicationUrl from environment variable: ", apiServerApplicationUrlStringEnvVar)
		configuration.ApiServer.ApplicationUrl = apiServerApplicationUrlStringEnvVar
	}

	if natsUrlStringEnvVar := os.Getenv("NATSUrl"); natsUrlStringEnvVar != "" {
		log.Println("Overriding NATSUrl from environment variable: ", natsUrlStringEnvVar)
		configuration.NATSUrl = natsUrlStringEnvVar
	}

	return configuration
}
