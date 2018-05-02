package libslack

import "github.com/Noah-Huppert/should-have-worked-from-home/msg"
import "github.com/nlopes/slack"
import "context"
import "fmt"
import "errors"

var NoSource error = errors.New("no Slack Source with id found")

// GetSource attempts to return a Source with the specified Slack ID. If no
// sources are found via the Slack API a NoSource error is returned.
func GetSource(ctx context.Context, api *slack.Client, id string) (*msg.Source,
	error) {

	// Attempt to find conversation with id
	conv, err := GetConversation(ctx, api, id)
	if err != nil {
		return nil, fmt.Errorf("error attempting to find conversation "+
			"with id: %s", err.Error())
	}

	if conv != nil {
		return conv, nil
	}

	// Attempt to find user with id
	user, err := GetUser(ctx, api, id)
	if err != nil {
		return nil, fmt.Errorf("error attempting to find user with "+
			"id: %s", err.Error())
	}

	if user != nil {
		return user, nil
	}

	// Not found
	return nil, NoSource
}
