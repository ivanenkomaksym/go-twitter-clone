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

	googleApplicationCredentialsEnvVar := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	log.Println("Read GOOGLE_APPLICATION_CREDENTIALS from environment variable: ", googleApplicationCredentialsEnvVar)

	// Set default google endpoints
	if authenticationEnableEnvVar := os.Getenv("AUTHENTICATION_ENABLE"); authenticationEnableEnvVar != "" {
		log.Println("Overriding AUTHENTICATION_ENABLE from environment variable: ", authenticationEnableEnvVar)
		configuration.Authentication.Enable, _ = strconv.ParseBool(authenticationEnableEnvVar)
	}

	if configuration.Authentication.Enable {
		configuration.Authentication.OAuth2.Endpoint = google.Endpoint
	}

	if modeEnvVar := os.Getenv("MODE"); modeEnvVar != "" {
		log.Println("Overriding MODE from environment variable: ", modeEnvVar)
		configuration.Mode, _ = ParseMode(modeEnvVar)
	}

	if projectIdEnvVar := os.Getenv("PROJECT_ID"); projectIdEnvVar != "" {
		log.Println("Overriding PROJECT_ID from environment variable: ", projectIdEnvVar)
		configuration.ProjectId = projectIdEnvVar
	}

	if clientIdEnvVar := os.Getenv("AUTHENTICATION_OAUTH2_CLIENTID"); clientIdEnvVar != "" {
		log.Println("Overriding AUTHENTICATION_OAUTH2_CLIENTID from environment variable: ", clientIdEnvVar)
		configuration.Authentication.OAuth2.ClientID = clientIdEnvVar
	}

	if clientSecretEnvVar := os.Getenv("AUTHENTICATION_OAUTH2_CLIENTSECRET"); clientSecretEnvVar != "" {
		log.Println("Overriding AUTHENTICATION_OAUTH2_CLIENTSECRET from environment variable: ", clientSecretEnvVar)
		configuration.Authentication.OAuth2.ClientSecret = clientSecretEnvVar
	}

	if redirectUrlEnvVar := os.Getenv("AUTHENTICATION_OAUTH2_REDIRECT_URI"); redirectUrlEnvVar != "" {
		log.Println("Overriding AUTHENTICATION_OAUTH2_REDIRECT_URI from environment variable: ", redirectUrlEnvVar)
		configuration.Authentication.OAuth2.RedirectURL = redirectUrlEnvVar
	}

	if tweetsStorageConnectionStringEnvVar := os.Getenv("TWEETSSTORAGE_CONNECTIONSTRING"); tweetsStorageConnectionStringEnvVar != "" {
		log.Println("Overriding TWEETSSTORAGE_CONNECTIONSTRING from environment variable: ", tweetsStorageConnectionStringEnvVar)
		configuration.TweetsStorage.ConnectionString = tweetsStorageConnectionStringEnvVar
	}

	if feedsStorageConnectionStringEnvVar := os.Getenv("FEEDSSTORAGE_CONNECTIONSTRING"); feedsStorageConnectionStringEnvVar != "" {
		log.Println("Overriding FEEDSSTORAGE_CONNECTIONSTRING from environment variable: ", feedsStorageConnectionStringEnvVar)
		configuration.FeedsStorage.ConnectionString = feedsStorageConnectionStringEnvVar
	}

	if apiServerApplicationUrlStringEnvVar := os.Getenv("APISERVER_APPLICATIONURL"); apiServerApplicationUrlStringEnvVar != "" {
		log.Println("Overriding APISERVER_APPLICATIONURL from environment variable: ", apiServerApplicationUrlStringEnvVar)
		configuration.ApiServer.ApplicationUrl = apiServerApplicationUrlStringEnvVar
	}

	if natsUrlStringEnvVar := os.Getenv("NATS_URL"); natsUrlStringEnvVar != "" {
		log.Println("Overriding NATS_URL from environment variable: ", natsUrlStringEnvVar)
		configuration.NATSUrl = natsUrlStringEnvVar
	}

	if redirectUriStringEnvVar := os.Getenv("REDIRECT_URI"); redirectUriStringEnvVar != "" {
		log.Println("Overriding REDIRECT_URI from environment variable: ", redirectUriStringEnvVar)
		configuration.RedirectURI = redirectUriStringEnvVar
	}

	if allowOriginStringEnvVar := os.Getenv("ALLOW_ORIGIN"); allowOriginStringEnvVar != "" {
		log.Println("Overriding ALLOW_ORIGIN from environment variable: ", allowOriginStringEnvVar)
		configuration.AllowOrigin = allowOriginStringEnvVar
	}

	configurationJson, err := json.MarshalIndent(configuration, "", "  ")
	if err != nil {
		log.Println("Error marshalling configuration to JSON: ", err)
	} else {
		log.Println("Configuration: ", string(configurationJson))
	}

	return configuration
}
