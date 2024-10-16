package config

import (
	_ "embed"
	"encoding/json"
	"log"
	"os"

	"golang.org/x/oauth2/google"
)

//go:embed appsettings.json
var appsettingsContent []byte

func ReadConfiguration() Configuration {
	configuration := Configuration{}
	json.Unmarshal(appsettingsContent, &configuration)

	// Set default google endpoints
	if configuration.Authentication.Enable {
		configuration.Authentication.OAuth2.Endpoint = google.Endpoint
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
