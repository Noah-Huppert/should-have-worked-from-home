package bot

import "github.com/nlopes/slack"
import "fmt"
import "context"

// GetChannels returns a mapping from Slack channel ids to human readable names
func GetChannels(ctx context.Context, api *slack.Client) (
	map[string]string, error) {

	// Get channels from API
	chans, err := api.GetChannelsContext(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("error getting channels: %s",
			err.Error())
	}

	// Map
	mapping := make(map[string]string)

	for _, c := range chans {
		mapping[c.ID] = c.Name
	}

	return mapping, nil
}
