package bot

import "github.com/Noah-Huppert/should-have-worked-from-home/msg"
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

	// Start Slack real time messenger listener
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	// Get mapping from source ids to names
	logger.Println("retrieving workspace metadata")

	sources, err := GetSources(ctx, api)
	if err != nil {
		errs <- fmt.Errorf("error getting sources mapping: %s",
			err.Error())
		errs <- nil
		return
	}

	for _, v := range sources {
		logger.Println(v)
	}

	logger.Println("started listening for Slack events")

	for {
		select {
		case <-ctx.Done():
			errs <- ctx.Err()
			errs <- nil
			return

		case msg := <-rtm.IncomingEvents:
			switch msg.Data.(type) {
			case *slack.ConnectedEvent:
				logger.Println("successfully connected to Slack API")

			case *slack.MessageEvent:
				msgEv, ok := msg.Data.(*slack.MessageEvent)
				if !ok {
					errs <- fmt.Errorf("error converting "+
						"message to MessageEvent "+
						"type: %#v", msg)
					continue
				}

				chanName, ok := sources[msgEv.Channel]
				if !ok {
					errs <- fmt.Errorf("could not find "+
						"message channel name, "+
						"channel id: %s, dropping "+
						"message", msgEv.Channel)
					continue
				}

				logger.Printf("received Slack message in "+
					"%s (chan: %s)\n",
					msgEv.Text, chanName)

			case *slack.RTMError:
				errs <- fmt.Errorf("slack RTM error: %#v",
					msg)

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
