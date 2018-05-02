package bot

import "github.com/nlopes/slack"
import "github.com/Noah-Huppert/should-have-worked-from-home/msg"
import "fmt"
import "context"

// GetSources returns a mapping from Slack ids to Source instances
func GetSources(ctx context.Context, api *slack.Client) (
	map[string]*msg.Source, error) {

	// Get users
	users, err := GetUsers(ctx, api)
	if err != nil {
		return nil, fmt.Errorf("error getting users: %s",
			err.Error())
	}

	/*
		// Get channels
			chans, err := GetChannels(ctx, api)
			if err != nil {
				return nil, fmt.Errorf("error getting channels: %s",
					err.Error())
			}
	*/

	// Get conversations
	convs, err := GetConversations(ctx, api)
	if err != nil {
		return nil, fmt.Errorf("error getting conversations: %s",
			err.Error())
	}

	// Map
	mapping := make(map[string]*msg.Source)

	for k, v := range users {
		mapping[k] = v
	}

	/*
		for k, v := range chans {
			s := msg.NewChannelSource(k, v)
			mapping[k] = s
		}
	*/

	for k, v := range convs {
		mapping[k] = v
	}

	return mapping, nil
}
