package bot

import "github.com/nlopes/slack"
import "github.com/Noah-Huppert/should-have-worked-from-home/msg"
import "fmt"
import "context"

// GetUsers returns a mapping from Slack user ids to Sources
func GetUsers(ctx context.Context, api *slack.Client) (
	map[string]*msg.Source, error) {

	// Get users from API
	users, err := api.GetUsersContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting users: %s",
			err.Error())
	}

	// Map
	mapping := make(map[string]*msg.Source)

	for _, u := range users {
		s := msg.NewSource(u.ID, u.Name, msg.User)
		mapping[u.ID] = s
	}

	return mapping, nil
}
