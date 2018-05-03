package msg

import "time"
import "fmt"

// Msg provides details about a Slack message which indicates a certain user
// should have worked from home today
type Msg struct {
	// In is where the message was sent from
	In *Source

	// Sender is the source of the message
	Sender *Source

	// SentAt is the time the message was sent
	SentAt *time.Time

	// Subject the user who should have worked from home today
	Subject *Source

	// Text is the message contents
	Text string
}

// String returns a text representation of a message
func (m Msg) String() string {
	return fmt.Sprintf("[In: %s\n"+
		"Sender: %s\n"+
		"SentAt: %s\n"+
		"Subject: %s\n"+
		"Text: %s]", m.In, m.Sender, m.SentAt, m.Subject, m.Text)
}

// New creates a new Msg
func NewMsg(in *Source, sender *Source, sentAt *time.Time, subject *Source,
	text string) *Msg {

	return &Msg{
		In:      in,
		Sender:  sender,
		SentAt:  sentAt,
		Subject: subject,
		Text:    text,
	}
}
