package bot

import "github.com/nlopes/slack"
import "github.com/Noah-Huppert/should-have-worked-from-home/msg"
import "fmt"
import "context"

// GetConversations returns a mapping between Slack conversation ids and
// sources
func GetConversations(ctx context.Context, api *slack.Client) (
	map[string]*msg.Source, error) {

	// Get conversations from API
	mapping := make(map[string]*msg.Source)

	firstCall := true
	nextCursor := ""

	for firstCall || len(nextCursor) != 0 {
		if firstCall {
			firstCall = false
		}

		// Make request
		params := &slack.GetConversationsParameters{
			Cursor: nextCursor,
		}

		convs, next, err := api.GetConversationsContext(ctx, params)
		nextCursor = next

		// Handle request
		if err != nil {
			return nil, fmt.Errorf("error getting conversations: %s",
				err.Error())
		}

		for _, c := range convs {
			s := msg.Source{
				ID:   c.ID,
				Name: c.Name,
			}

			if c.IsChannel || c.IsGroup {
				s.Type = msg.Channel
			} else if c.IsIM {
				s.Type = msg.IM
			} else {
				return nil, fmt.Errorf("error determining "+
					"conversation type: %#v", c)
			}

			mapping[c.ID] = &s
		}
	}

	return mapping, nil

}
