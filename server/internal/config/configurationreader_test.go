package config_test

import (
	"os"
	"reflect"
	"strconv"
	"testing"
	config "twitter-clone/internal/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func TestReadConfiguration_DefaultsFromEmbeddedFile(t *testing.T) {
	// Clear environment variables to use only appsettings.json defaults
	os.Clearenv()

	expectedMode, _ := config.ParseMode("Cloud")
	expected := config.Configuration{
		Mode:      expectedMode,
		ProjectId: "{your-project-id}",
		ApiServer: config.ApiServer{
			ApplicationUrl: "localhost:8016",
		},
		TweetsStorage: config.TweetsStorage{
			ConnectionString: "myuser:mypassword@tcp(127.0.0.1:3306)",
			DatabaseName:     "TweetsDb",
		},
		FeedsStorage: config.FeedsStorage{
			ConnectionString: "mongodb://localhost:27017",
			DatabaseName:     "FeedsDb",
			CollectionName:   "Feeds",
		},
		RedirectURI: "http://localhost:3000/callback",
		AllowOrigin: "http://localhost:3000",
		Authentication: config.Authentication{
			Enable: true,
			OAuth2: oauth2.Config{
				RedirectURL:  "http://localhost:8016/auth/google/callback",
				ClientID:     "{CLIENTID}.apps.googleusercontent.com",
				ClientSecret: "{SECRET}",
				Scopes: []string{
					"https://www.googleapis.com/auth/userinfo.email",
					"https://www.googleapis.com/auth/userinfo.profile",
				},
				Endpoint: google.Endpoint,
			},
		},
	}

	config := config.ReadConfiguration()

	if !reflect.DeepEqual(config, expected) {
		t.Errorf("Expected %+v config to match defaults from appsettings.json, got: %+v", expected, config)
	}
}

func TestReadConfiguration_WithEnvVars(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("AUTHENTICATION_ENABLE", "true")
	os.Setenv("MODE", "development")
	os.Setenv("PROJECT_ID", "test-project-id")
	os.Setenv("AUTHENTICATION_OAUTH2_CLIENTID", "test-client-id")
	os.Setenv("AUTHENTICATION_OAUTH2_CLIENTSECRET", "test-client-secret")
	os.Setenv("TWEETSSTORAGE_CONNECTIONSTRING", "test-tweets-connection")
	os.Setenv("FEEDSSTORAGE_CONNECTIONSTRING", "test-feeds-connection")
	os.Setenv("APISERVER_APPLICATIONURL", "http://localhost:8080")
	os.Setenv("NATS_URL", "nats://localhost:4222")

	defer func() {
		// Clean up environment variables after the test
		os.Clearenv()
	}()

	// Call the function to test
	configuration := config.ReadConfiguration()

	// Verify that environment variables override the default configuration
	if !configuration.Authentication.Enable {
		t.Errorf("Expected Authentication.Enable to be true, got false")
	}

	expectedMode, _ := config.ParseMode("development")
	if configuration.Mode != expectedMode {
		t.Errorf("Expected Mode to be 'development', got %v", configuration.Mode)
	}

	if configuration.ProjectId != "test-project-id" {
		t.Errorf("Expected ProjectId to be 'test-project-id', got %v", configuration.ProjectId)
	}

	if configuration.Authentication.OAuth2.ClientID != "test-client-id" {
		t.Errorf("Expected OAuth2.ClientID to be 'test-client-id', got %v", configuration.Authentication.OAuth2.ClientID)
	}

	if configuration.Authentication.OAuth2.ClientSecret != "test-client-secret" {
		t.Errorf("Expected OAuth2.ClientSecret to be 'test-client-secret', got %v", configuration.Authentication.OAuth2.ClientSecret)
	}

	if configuration.TweetsStorage.ConnectionString != "test-tweets-connection" {
		t.Errorf("Expected TweetsStorage.ConnectionString to be 'test-tweets-connection', got %v", configuration.TweetsStorage.ConnectionString)
	}

	if configuration.FeedsStorage.ConnectionString != "test-feeds-connection" {
		t.Errorf("Expected FeedsStorage.ConnectionString to be 'test-feeds-connection', got %v", configuration.FeedsStorage.ConnectionString)
	}

	if configuration.ApiServer.ApplicationUrl != "http://localhost:8080" {
		t.Errorf("Expected ApiServer.ApplicationUrl to be 'http://localhost:8080', got %v", configuration.ApiServer.ApplicationUrl)
	}

	if configuration.NATSUrl != "nats://localhost:4222" {
		t.Errorf("Expected NATSUrl to be 'nats://localhost:4222', got %v", configuration.NATSUrl)
	}
}

func TestReadConfiguration_EnvironmentVariableOverrides(t *testing.T) {
	os.Setenv("AUTHENTICATION_ENABLE", "false")
	defer os.Clearenv()

	// Check if the `Authentication.Enable` gets overridden correctly by the environment variable
	config := config.ReadConfiguration()

	expected, _ := strconv.ParseBool(os.Getenv("AUTHENTICATION_ENABLE"))
	if config.Authentication.Enable != expected {
		t.Errorf("Expected Authentication.Enable to be %v, got %v", expected, config.Authentication.Enable)
	}
}
