package bot

import "github.com/Noah-Huppert/should-have-worked-from-home/msg"
import "github.com/Noah-Huppert/should-have-worked-from-home/libslack"
import "github.com/Noah-Huppert/should-have-worked-from-home/config"
import "os"
import "log"
import "fmt"
import "context"
import "github.com/nlopes/slack"

// Listen watches a Slack channel for a "I should have worked from home today"
// message and signals a channel when one is received.
//
// Additionally an error channel is returned. The listener will attempt to
// recover from errors. However if `nil` is sent over the errors channel then
// the listener was unable to recover, and exited.
func Listen(ctx context.Context, cfg *config.Config) (<-chan *msg.Msg,
	<-chan error) {

	logger := log.New(os.Stdout, "listen: ", 0)

	// Setup Slack
	api := slack.New(cfg.SlackToken)
	slack.SetLogger(logger)

	// Handle messages
	msgs := make(chan *msg.Msg)
	errs := make(chan error)

	go handleSlackEvents(ctx, cfg, logger, api, msgs, errs)

	return msgs, errs
}

// handleSlackEvents performs the correct action for each received Slack event
func handleSlackEvents(ctx context.Context, cfg *config.Config,
	logger *log.Logger, api *slack.Client, msgs chan *msg.Msg,
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
	msgs chan *msg.Msg, errs chan error) {

	msg, ok := event.Data.(*slack.MessageEvent)
	if !ok {
		errs <- fmt.Errorf("error converting message to MessageEvent "+
			"type: %#v", msg)
		return
	}

	/*
		chanName, ok := sources[msgEv.Channel]
		if !ok {
			errs <- fmt.Errorf("could not find "+
				"message channel name, "+
				"channel id: %s, dropping "+
				"message", msgEv.Channel)
			continue
		}
	*/

	// Lazy load message source
	sourceId := msg.Channel
	source, ok := sources[sourceId]

	if !ok {
		s, err := libslack.GetConversation(ctx, api, sources, sourceId)
		if err != nil {
			errs <- fmt.Errorf("error finding message source: %s",
				err.Error())
			return
		}

		sources[sourceId] = s
		source = s
	}

	// TODO: Load sender of message
	// TODO: Construct Msg instance
	// TODO: Record time message received
	// TODO: Send Msg instance in channel

	//msg := msg.NewMsg()
	logger.Printf("received Slack message: %s, from: %s\n", msg.Text,
		source)
}
