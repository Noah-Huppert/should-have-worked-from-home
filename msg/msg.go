package msg

import "time"
import "fmt"

// Msg provides details on a received Slack message
type Msg struct {
	// Sender is the name of the user who sent the message
	Sender string

	// SentAt is the time the message was sent
	SentAt time.Time

	// Text is the message contents
	Text string
}

// String returns a text representation of a message
func (m Msg) String() string {
	return fmt.Sprintf("[%s] %s: %s", m.SentAt, m.Sender, m.Text)
}

// New creates a new Msg
func New(Sender string, SentAt time.Time, Text string) *Msg {
	return &Msg{
		Sender: Sender,
		SentAt: SentAt,
		Text:   Text,
	}
}
