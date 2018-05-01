package bot

import "github.com/nlopes/slack"
import "fmt"
import "context"

// GetUsers returns a mapping from Slack user ids to human readable names
func GetUsers(ctx context.Context, api *slack.Client) (
	map[string]string, error) {

	// Get users from API
	users, err := api.GetUsersContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting users: %s",
			err.Error())
	}

	// Map
	mapping := make(map[string]string)

	for _, u := range users {
		mapping[u.ID] = u.Name
	}

	return mapping, nil
}
