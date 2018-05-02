package libslack

import "github.com/Noah-Huppert/should-have-worked-from-home/msg"
import "github.com/nlopes/slack"
import "context"
import "fmt"

// GetConversation returns a Source for a conversation id
func GetConversation(ctx context.Context, api *slack.Client, id string) (
	*msg.Source, error) {

	// Get conversation from API
	conv, err := api.GetConversationInfoContext(ctx, id, true)
	if err != nil {
		return nil, fmt.Errorf("error retrieving conversation info: "+
			" id: %s, err: %s", id, err.Error())
	}

	// Create Source
	s := &msg.Source{
		ID:   conv.ID,
		Name: conv.Name,
	}

	if conv.IsChannel || conv.IsGroup {
		s.Type = msg.Channel
	} else if conv.IsIM {
		s.Type = msg.IM
	} else {
		return nil, fmt.Errorf("error determining "+
			"conversation type: %#v", conv)
	}

	return s, nil
}
