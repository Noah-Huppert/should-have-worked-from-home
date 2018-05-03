package libslack

import "github.com/Noah-Huppert/should-have-worked-from-home/msg"
import "github.com/nlopes/slack"
import "context"
import "fmt"

// GetUser returns a Source representing a Slack user
func GetUser(ctx context.Context, api *slack.Client,
	sources map[string]*msg.Source, id string) (*msg.Source, error) {

	if user, ok := sources[id]; ok {
		return user, nil
	}

	// Get user from API
	user, err := api.GetUserInfoContext(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user info: %s",
			err.Error())
	}

	// Make source
	s := msg.NewSource(user.ID, user.Name, msg.User)

	return s, nil
}
