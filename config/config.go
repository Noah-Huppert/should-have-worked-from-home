package config

import "os"
import "fmt"

// EnvKeySlackToken is the environment variable key the Slack token will be
// provided under
const EnvKeySlackToken string = "SLACK_TOKEN"

type Config struct {
	// SlackToken is the Slack API token
	SlackToken string
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

	// Make instance
	c := Config{
		SlackToken: slackToken,
	}

	return &c, nil
}
