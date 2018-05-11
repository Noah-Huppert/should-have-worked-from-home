package config

import "os"
import "fmt"
import "errors"

// EnvKeySlackToken is the environment variable key the Slack token will be
// provided by
const EnvKeySlackToken string = "SLACK_TOKEN"

// EnvKeyGAPIAccessToken is the environment variable key the GAPI access token
// will be provided by
const EnvKeyGAPIAccessToken string = "GAPI_ACCESS_TOKEN"

// EnvKeyGAPIPrivateKeyPath is the environment variable key the path to the
// GAPI private key file will be provided by
const EnvKeyGAPIPrivateKeyPath string = "GAPI_PRIVATE_KEY_PATH"

// EnvKeySpreadsheetID is the environment variable key in which values for
// Config.SpreadsheetID will be provided
const EnvKeySpreadsheetID string = "SPREADSHEET_ID"

// EnvKeySpreadsheetPageName is the environment variable key which values for
// Config.SpreadsheetPageName will be provided
const EnvKeySpreadsheetPageName string = "SPREADSHEET_PAGE_NAME"

// DefaultGAPIKeyPath is default location of the GAPI key file if none is
// provided
const DefaultGAPIPrivateKeyPath string = "gapi_private_key.json"

// ErrNoGAPIAccessToken indicates a GAPI access token was not provided
var ErrNoGAPIAccessToken error = errors.New("no GAPI access token provided")

type Config struct {
	// SlackToken is the Slack API token
	SlackToken string

	// GAPIAccessToken is the access token used to authenticate with the
	// GAPI to get and modify a user's spreadsheet
	GAPIAccessToken string

	// GAPIPrivateKeyPath is the path to the Google API private key file.
	// It defaults to DefaultGAPIPrivateKeyPath if not specified by an
	// environment variable
	GAPIPrivateKeyPath string

	// SpreadsheetID is the id of the Google spreadsheet which will hold
	// message information
	SpreadsheetID string

	// SpreadsheetPageName is the name of the Google Spreadsheet page
	// which "I should have worked from home today" message will be
	// recorded in
	SpreadsheetPageName string
}

// New constructs a new Config from environment variables. An error is returned
// if any required environment variables are not set.
//
// If ErrNoGAPIAccessToken is provided the user should be directed to retrieve a
// GAPI access token via the URL returned by gapi.GetTokenURL().
func New() (*Config, error) {
	// Get SlackToken
	slackToken, ok := os.LookupEnv(EnvKeySlackToken)
	if !ok {
		return nil, fmt.Errorf("environment variable \"%s\" must be "+
			"provided", EnvKeySlackToken)
	}

	// Get GAPIPrivateKeyPath if set
	gapiPrivateKeyPath, ok := os.LookupEnv(EnvKeyGAPIPrivateKeyPath)
	if !ok {
		gapiPrivateKeyPath = DefaultGAPIPrivateKeyPath
	}

	// Get GAPIAccessToken
	gapiAccessToken, ok := os.LookupEnv(EnvKeyGAPIAccessToken)
	var gapiAccessTokenErr error = nil
	if !ok {
		gapiAccessToken = ""
		gapiAccessTokenErr = ErrNoGAPIAccessToken
	}

	// Get SpreadsheetID
	spreadsheetID, ok := os.LookupEnv(EnvKeySpreadsheetID)
	if !ok {
		return nil, fmt.Errorf("environment variable \"%s\" must be "+
			"provided", EnvKeySpreadsheetID)
	}

	// Get SpreadsheetPageName
	spreadsheetPageName, ok := os.LookupEnv(EnvKeySpreadsheetPageName)
	if !ok {
		return nil, fmt.Errorf("environment variable \"%s\" must be "+
			"provided", EnvKeySpreadsheetPageName)
	}

	// Make instance
	c := &Config{
		SlackToken:          slackToken,
		GAPIAccessToken:     gapiAccessToken,
		GAPIPrivateKeyPath:  gapiPrivateKeyPath,
		SpreadsheetID:       spreadsheetID,
		SpreadsheetPageName: spreadsheetPageName,
	}

	// If error retrieving GAPI access token
	if gapiAccessTokenErr != nil {
		return c, gapiAccessTokenErr
	}

	return c, nil
}
