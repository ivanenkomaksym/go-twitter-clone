package config

import (
	_ "embed"
	"encoding/json"
	"os"
)

//go:embed appsettings.json
var appsettingsContent []byte

func ReadConfiguration() Configuration {
	configuration := Configuration{}
	json.Unmarshal(appsettingsContent, &configuration)

	if tweetsStorageConnectionStringEnvVar := os.Getenv("TweetsStorage__ConnectionString"); tweetsStorageConnectionStringEnvVar != "" {
		configuration.TweetsStorage.ConnectionString = tweetsStorageConnectionStringEnvVar
	}

	if feedsStorageConnectionStringEnvVar := os.Getenv("FeedsStorage__ConnectionString"); feedsStorageConnectionStringEnvVar != "" {
		configuration.FeedsStorage.ConnectionString = feedsStorageConnectionStringEnvVar
	}

	if apiServerApplicationUrlStringEnvVar := os.Getenv("ApiServer__ApplicationUrl"); apiServerApplicationUrlStringEnvVar != "" {
		configuration.ApiServer.ApplicationUrl = apiServerApplicationUrlStringEnvVar
	}

	return configuration
}
