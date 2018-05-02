package msg

import "time"
import "fmt"

// Msg provides details on a received Slack message
type Msg struct {
	// In is where the message was sent from
	In *Source

	// Sender is the source of the message
	Sender *Source

	// SentAt is the time the message was sent
	SentAt time.Time

	// Text is the message contents
	Text string
}

// String returns a text representation of a message
func (m Msg) String() string {
	return fmt.Sprintf("[%s - %s] %s: %s", m.In, m.SentAt, m.Sender, m.Text)
}

// New creates a new Msg
func New(in *Source, Sender *Source, SentAt time.Time, Text string) *Msg {
	return &Msg{
		In:     in,
		Sender: Sender,
		SentAt: SentAt,
		Text:   Text,
	}
}
