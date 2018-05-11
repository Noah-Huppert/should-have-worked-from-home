package config

import "os"
import "fmt"
import "errors"

// EnvKeySlackToken is the environment variable key the Slack token will be
// provided by
const EnvKeySlackToken string = "SLACK_TOKEN"

// EnvKeyGAPIOAuthToken is the environment variable key the GAPI OAuth token
// will be provided by
const EnvKeyGAPIOAuthToken string = "GAPI_OAUTH_TOKEN"

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

// ErrNoGAPIOAuthToken indicates a GAPI OAuth token was not provided
var ErrNoGAPIOAuthToken error = errors.New("no GAPI OAuth token provided")

type Config struct {
	// SlackToken is the Slack API token
	SlackToken string

	// GAPIOAuthToken is the OAuth token used to authenticate with the
	// GAPI to access the user's spreadsheet
	GAPIOAuthToken string

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
// If ErrNoGAPIOAuthToken is provided the user should be directed to retrieve a
// GAPI OAuth token via the token URL.
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

	// Get GAPIOAuthToken
	gapiOAuthToken, ok := os.LookupEnv(EnvKeyGAPIOAuthToken)
	var gapiOAuthTokenErr error = nil
	if !ok {
		gapiOAuthToken = ""
		gapiOAuthTokenErr = ErrNoGAPIOAuthToken
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
		GAPIOAuthToken:      gapiOAuthToken,
		GAPIPrivateKeyPath:  gapiPrivateKeyPath,
		SpreadsheetID:       spreadsheetID,
		SpreadsheetPageName: spreadsheetPageName,
	}

	// If error retrieving GAPI OAuth token
	if gapiOAuthTokenErr != nil {
		return c, gapiOAuthTokenErr
	}

	return c, nil
}
