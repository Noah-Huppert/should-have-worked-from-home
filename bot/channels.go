package bot

import "github.com/nlopes/slack"
import "fmt"
import "context"

// GetChannelId returns the Slack channel id for the provided channel name.
func GetChannelId(ctx context.Context, api *slack.Client, name string) (string,
	error) {

	// Get channels from API
	chans, err := api.GetChannelsContext(ctx, false)
	if err != nil {
		return "", fmt.Errorf("error getting channel \"%s\" id: %s",
			name, err.Error())
	}

	// Search
	for _, c := range chans {
		if c.Name == name {
			return c.ID, nil
		}
	}

	// None found
	return "", fmt.Errorf("no channels with name \"%s\" found", name)
}
