package libslack

import "github.com/nlopes/slack"
import "context"
import "fmt"

// GetConversationMembers returns the users in a conversation
func GetConversationMembers(ctx context.Context, api *slack.Client, id string) (
	[]string, error) {

	members := []string{}
	cursor := ""
	firstCall := true

	for firstCall || len(cursor) != 0 {
		if firstCall {
			firstCall = false
		}

		newMems, newCursor, err := getMembers(ctx, api, members, id, cursor)

		members = newMems
		cursor = newCursor

		if err != nil {
			return nil, err
		}
	}

	return members, nil
}

// getMembers is used by GetConversationMembers to make the
// GetUsersInConversation Slack API call. The new list of members and the
// next request cursor are returned.
func getMembers(ctx context.Context, api *slack.Client, members []string,
	id string, cursor string) ([]string, string, error) {

	params := &slack.GetUsersInConversationParameters{
		ChannelID: id,
		Cursor:    cursor,
	}
	newMembers, nextCursor, err := api.GetUsersInConversationContext(ctx, params)

	if err != nil {
		return nil, "", fmt.Errorf("error retrieving users in conversation, "+
			"conversation id: %s, cursor: %s, err: %s", id, cursor,
			err.Error())
	}

	for _, member := range newMembers {
		members = append(members, member)
	}

	return members, nextCursor, nil
}
