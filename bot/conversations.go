package bot

import "github.com/nlopes/slack"
import "fmt"
import "context"

// GetConversations returns a mapping between Slack conversation ids and human
// readable names
func GetConversations(ctx context.Context, api *slack.Client) (
	map[string]string, error) {

	// Get conversations from API
	mapping := make(map[string]string)

	firstCall := true
	nextCursor := ""

	for firstCall || len(nextCursor) != 0 {
		if firstCall {
			firstCall = false
		}

		params := &slack.GetConversationsParameters{
			Cursor: nextCursor,
		}

		convs, next, err := api.GetConversationsContext(ctx, params)
		nextCursor = next
		if err != nil {
			return nil, fmt.Errorf("error getting conversations: %s",
				err.Error())
		}

		for _, c := range convs {
			fmt.Println(c.ID, c.Name)
			mapping[c.ID] = c.Name
		}
	}
	// TODO: Convert to GetIMs
	// TODO: Make GetGroups method

	return mapping, nil

}
