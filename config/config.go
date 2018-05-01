package config

import "os"
import "fmt"

// EnvKeySlackToken is the environment variable key the Slack token will be
// provided under
const EnvKeySlackToken string = "SLACK_TOKEN"

// EnvKeySlackChannel is the environment variable key the Slack channel will be
// provided by
const EnvKeySlackChannel string = "SLACK_CHANNEL"

// DefaultSlackChannel the default Slack channel value
const DefaultSlackChannel string = "should-have-worked-from-home"

type Config struct {
	// SlackToken is the Slack API token
	SlackToken string

	// SlackChannel is the Slack channel to listen for "I should have
	// worked from home today" messages in
	SlackChannel string
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

	// Get SlackChannel
	slackChan, ok := os.LookupEnv(EnvKeySlackChannel)
	if !ok {
		slackChan = DefaultSlackChannel
	}

	// Make instance
	c := Config{
		SlackToken:   slackToken,
		SlackChannel: slackChan,
	}

	return &c, nil
}
