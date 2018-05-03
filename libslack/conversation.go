package libslack

import "github.com/Noah-Huppert/should-have-worked-from-home/msg"
import "github.com/nlopes/slack"
import "context"
import "fmt"
import "strings"

// GetConversation returns a Source for a conversation id
func GetConversation(ctx context.Context, api *slack.Client,
	sources map[string]*msg.Source, id string) (*msg.Source, error) {

	// If in sources cache
	if conv, ok := sources[id]; ok {
		return conv, nil
	}

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

	// -- -- Name
	if len(s.Name) == 0 {
		// If conversation doesn't have name, make one out of member
		// names
		name, err := makeConversationName(ctx, api, sources, s.ID)
		if err != nil {
			return nil, fmt.Errorf("error making name for "+
				"conversation: %s", err.Error())
		}

		s.Name = name
	}

	// -- -- Type
	if conv.IsChannel || conv.IsGroup {
		s.Type = msg.Channel
	} else if conv.IsIM {
		s.Type = msg.IM
	} else {
		return nil, fmt.Errorf("error determining "+
			"conversation type: %#v", conv)
	}

	// Add to sources cache
	sources[s.ID] = s

	return s, nil
}

// makeConversationName creates a name for a conversation out of the members
// of the conversation
func makeConversationName(ctx context.Context, api *slack.Client,
	sources map[string]*msg.Source, id string) (string, error) {

	// Get member ids
	members, err := GetConversationMembers(ctx, api, id)
	if err != nil {
		return "", fmt.Errorf("error getting conversation "+
			"members: %s", err.Error())
	}

	// Get member names
	memberNames := []string{}
	for _, memberId := range members {
		user, err := GetUser(ctx, api, sources, memberId)
		if err != nil {
			return "", fmt.Errorf("error retrieving "+
				"conversation member, id: %s, err: %s",
				memberId, err.Error())
		}

		memberNames = append(memberNames, user.Name)
	}

	if len(memberNames) == 0 {
		return "", fmt.Errorf("no member names retrieved for "+
			"conversation, id: %#v", id)
	}

	return strings.Join(memberNames, "--"), nil
}
