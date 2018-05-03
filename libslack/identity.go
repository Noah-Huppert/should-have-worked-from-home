package libslack

import "github.com/Noah-Huppert/should-have-worked-from-home/msg"
import "github.com/nlopes/slack"
import "context"
import "fmt"

// GetIdentity returns the Slack bot's own information
func GetIdentity(ctx context.Context, api *slack.Client,
	sources map[string]*msg.Source) (*msg.Source, error) {

	// Check sources cache
	if s, ok := sources[msg.SourcesBotSelfKey]; ok {
		return s, nil
	}

	// Get bot identity from api
	ident, err := api.GetBotInfoContext(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("error retrieving bot identity: %s",
			err.Error())
	}

	s := msg.NewSource(ident.ID, ident.Name, msg.User)

	// Cache identity
	sources[msg.SourcesBotSelfKey] = s

	return s, nil
}
