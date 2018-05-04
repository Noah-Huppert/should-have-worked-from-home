package libslack

import "github.com/Noah-Huppert/should-have-worked-from-home/msg"
import "github.com/nlopes/slack"
import "context"
import "fmt"
import "strings"
import "regexp"

// MentionExp matches the Slack @ mention API format (<@ID>).
// The first capture group will be the mentioned user's ID.
var MentionExp *regexp.Regexp = regexp.MustCompile("<@([A-Z0-9]*)>")

// TransformMentions converts Slack @ mentions from API format (<@ID>) to
// user readable format (@Username).
// The transformed version of the message text provided will be returned.
func TransformMentions(ctx context.Context, api *slack.Client,
	sources map[string]*msg.Source, msgText string) (string, error) {

	// Find all API formatted @ mentions
	matches := MentionExp.FindAllStringSubmatch(msgText, -1)

	// If no @ mentions
	if len(matches) == 0 {
		return msgText, nil
	}

	// Extract @ mention ids
	ids := []string{}
	for _, match := range matches {
		ids = append(ids, match[1])
	}

	// For each @ mention id
	for _, id := range ids {
		// Get user
		user, err := GetUser(ctx, api, sources, id)
		if err != nil {
			return "", fmt.Errorf("error retrieving user for @ "+
				"reference with id: %s, error: %s", id,
				err.Error())
		}

		// Transform
		msgText = strings.Replace(msgText,
			fmt.Sprintf("<@%s>", id),
			fmt.Sprintf("@%s", user.Name), -1)
	}

	return msgText, nil
}
