package bot

import "github.com/Noah-Huppert/should-have-worked-from-home/msg"
import "github.com/Noah-Huppert/should-have-worked-from-home/libslack"
import "github.com/Noah-Huppert/should-have-worked-from-home/config"
import "os"
import "log"
import "fmt"
import "context"
import "github.com/nlopes/slack"
import "regexp"
import "strings"

// MsgSubjectExp extract the subject of a "should have worked from home" message
// First capture group will either be "I" or a @ mention of a user
// Second capture group will hold the rest of a message which may contain a
//	"b/c reason I should have worked fro home today" piece of message
var MsgSubjectExp *regexp.Regexp = regexp.MustCompile(".*([iI]|<@[A-Z0-9]*>)" +
	".*[wW]orked.*[fF]rom.*[hH]ome.*")

// MsgReasonExp extract the reason for a "should have worked from home" message
// The first capture group is the reason
var MsgReasonExp *regexp.Regexp = regexp.MustCompile(".*(?:b\\/{0,1}c|because" +
	") +(.*)")

// Listen watches a Slack channel for a "I should have worked from home today"
// message and signals a channel when one is received.
//
// Additionally an error channel is returned. The listener will attempt to
// recover from errors. However if `nil` is sent over the errors channel then
// the listener was unable to recover, and exited.
func Listen(ctx context.Context, cfg *config.Config) (<-chan *msg.TargetMsg,
	<-chan error) {

	logger := log.New(os.Stdout, "listen: ", 0)

	// Setup Slack
	api := slack.New(cfg.SlackToken)
	slack.SetLogger(logger)

	// Handle messages
	msgs := make(chan *msg.TargetMsg)
	errs := make(chan error)

	go handleSlackEvents(ctx, cfg, logger, api, msgs, errs)

	return msgs, errs
}

// handleSlackEvents performs the correct action for each received Slack event
func handleSlackEvents(ctx context.Context, cfg *config.Config,
	logger *log.Logger, api *slack.Client, msgs chan *msg.TargetMsg,
	errs chan error) {

	// sources holds a mapping from Slack ids to Source instances
	sources := make(map[string]*msg.Source)

	// Start Slack real time messenger listener
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	logger.Println("started listening for Slack events")

	for {
		select {
		case <-ctx.Done():
			errs <- ctx.Err()
			errs <- nil
			return

		case event := <-rtm.IncomingEvents:
			switch event.Data.(type) {
			case *slack.ConnectedEvent:
				logger.Println("successfully connected to Slack API")

			case *slack.MessageEvent:
				handleMessage(ctx, api, logger, sources, event,
					msgs, errs)

			case *slack.RTMError:
				errs <- fmt.Errorf("slack RTM error: %#v",
					event)

			case *slack.InvalidAuthEvent:
				errs <- fmt.Errorf("invalid Slack " +
					"authentication")

				// Exit listener
				errs <- nil
				return

			default:
				continue
			}
		}
	}
}

// handleMessage is invoked when a Slack message event is received
func handleMessage(ctx context.Context, api *slack.Client, logger *log.Logger,
	sources map[string]*msg.Source, event slack.RTMEvent,
	msgs chan *msg.TargetMsg, errs chan error) {

	// Convert to MessageEvent
	msgEvent, ok := event.Data.(*slack.MessageEvent)
	if !ok {
		errs <- fmt.Errorf("error converting message to MessageEvent "+
			"type, skipping: %#v", msgEvent)
		return
	}

	// Load message source
	sourceId := msgEvent.Channel

	source, err := libslack.GetConversation(ctx, api, sources, sourceId)
	if err != nil {
		errs <- fmt.Errorf("error finding message source, skipping, "+
			"message: %#v, error: %s", msgEvent, err.Error())
		return
	}

	// Get message sender
	sender, err := libslack.GetUser(ctx, api, sources, msgEvent.User)
	if err != nil {
		errs <- fmt.Errorf("error getting message sender, skipping, "+
			"message: %#v, error: %s", msgEvent, err.Error())
	}

	// Get time message was sent
	sentAt, err := libslack.ConvertTimestamp(msgEvent.Timestamp)
	if err != nil {
		errs <- fmt.Errorf("error parsing message timestamp, "+
			"skipping, message: %#v, error: %s", msgEvent, err.Error())
		return
	}

	// Remove @ mentions of bot from message
	text, err := removeSelfMentions(ctx, api, sources, msgEvent.Text)
	if err != nil {
		errs <- fmt.Errorf("error processing message, skipping, "+
			"message: %#v, error: %s", msgEvent, err.Error())
		return
	}

	// Get subject of message
	subject, err := getMessageSubject(ctx, api, sources, msgEvent, text)
	if err != nil {
		errs <- fmt.Errorf("error determining message subject, "+
			"skipping, message: %#v, error: %s", msgEvent, err.Error())
		return
	}

	// If "should have worked from home" style message
	if subject != nil {
		// Transform message text @ mentions to be human readable
		transformedTxt, err := libslack.TransformMentions(ctx, api,
			sources, msgEvent.Text)
		if err != nil {
			errs <- fmt.Errorf("error transforming message @ "+
				"mentions to human readable format, "+
				"skipping..., message: %#v, error: %s",
				msgEvent, err.Error())
			return
		}

		// Get reason from message
		reason := getMessageReason(transformedTxt)

		// Send to channel
		m := msg.NewTargetMsg(source, sender, sentAt, subject, reason,
			transformedTxt)
		msgs <- m
	}
}

// getMessageReason determines the reason a "should have worked from home
// today" was sent.
func getMessageReason(msgText string) string {
	// Match
	matches := MsgReasonExp.FindStringSubmatch(msgText)

	// If no reason in message
	if len(matches) == 0 {
		return ""
	}

	// Get reason match
	return matches[1]
}

// getMessageSubject determines who the message is referring to. Returns a nil
// source if the message is not a "should have worked from home today" style
// message.
func getMessageSubject(ctx context.Context, api *slack.Client,
	sources map[string]*msg.Source, msgEvent *slack.MessageEvent,
	msgText string) (*msg.Source, error) {

	// Match
	matches := MsgSubjectExp.FindStringSubmatch(msgText)

	// Determine if "should have worked from home today" message
	if len(matches) == 0 {
		// Not in "should have worked from home today" form
		return nil, nil
	}

	// Get match subject
	subjText := matches[1]

	// Get user message is referencing
	var id string

	// If message is referencing sender of message
	if subjText == "i" || subjText == "I" {
		id = msgEvent.User
	} else { // Message is referencing other user
		// Extract id from @ mention format
		// <@ID>
		id = subjText[2 : len(subjText)-1]
	}

	user, err := libslack.GetUser(ctx, api, sources, id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user referenced "+
			"in message: %s", err.Error())
	}

	return user, nil

}

// removeSelfMentions strips a message of an @ mentions of the bot user
func removeSelfMentions(ctx context.Context, api *slack.Client,
	sources map[string]*msg.Source, msgText string) (string, error) {

	// Get bot identity
	ident, err := libslack.GetIdentity(ctx, api, sources)
	if err != nil {
		return "", fmt.Errorf("error stripping mentions of bot out "+
			"of message: %s", err.Error())
	}

	// Replace
	mentionStr := fmt.Sprintf("<@%s>", ident.ID)
	return strings.Replace(msgText, mentionStr, "", -1), nil
}
