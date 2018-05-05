package config

import "os"
import "fmt"

// EnvKeySlackToken is the environment variable key the Slack token will be
// provided by
const EnvKeySlackToken string = "SLACK_TOKEN"

// EnvKeyGAPIKeyPath is the environment variable key the path to the GAPI
// key file will be provided by
const EnvKeyGAPIKeyPath string = "GAPI_KEY_PATH"

// EnvKeySpreadsheetStore is the environment variable key which the id of the
// spreadsheet information will be stored in will be provided by
const EnvKeySpreadsheetStore string = "SPREADSHEET_STORE"

// DefaultGAPIKeyPath is default location of the GAPI key file if none is
// provided
const DefaultGAPIKeyPath string = "gapi_key.json"

type Config struct {
	// SlackToken is the Slack API token
	SlackToken string

	// GAPIKeyPath is the path to the Google API key file. It defaults to
	// DefaultGAPIKeyPath if not specified by an environment variable
	GAPIKeyPath string

	// SpreadsheetStore is the id of the Google spreadsheet which "I
	// should have worked from home today" messages will be recorded in
	SpreadsheetStore string
}

// New constructs a new Config from environment variables. An error is returned
// if any required environment variables are not set.
func New() (*Config, error) {
	// Get SlackToken
	slackToken, ok := os.LookupEnv(EnvKeySlackToken)
	if !ok {
		return nil, fmt.Errorf("environment variable \"%s\" must be "+
			"provided", EnvKeySlackToken)
	}

	// Get GAPIKeyPath if set
	gapiKeyPath, ok := os.LookupEnv(EnvKeyGAPIKeyPath)
	if !ok {
		gapiKeyPath = DefaultGAPIKeyPath
	}

	// Get SpreadsheetStore
	spreadsheetStore, ok := os.LookupEnv(EnvKeySpreadsheetStore)
	if !ok {
		return nil, fmt.Errorf("environment variable \"%s\" must be "+
			"provided", EnvKeySpreadsheetStore)
	}

	// Make instance
	c := Config{
		SlackToken:       slackToken,
		GAPIKeyPath:      gapiKeyPath,
		SpreadsheetStore: spreadsheetStore,
	}

	return &c, nil
}
