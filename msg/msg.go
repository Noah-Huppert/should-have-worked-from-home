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

	// Reason is the explanation for why a should have worked from home
	// today message was sent
	Reason string

	// Text is the message contents
	Text string
}

// String returns a text representation of a message
func (m Msg) String() string {
	sentAtStr := fmt.Sprintf("%d/%d/%d %d:%d", m.SentAt.Month(),
		m.SentAt.Day(), m.SentAt.Year(), m.SentAt.Hour(), m.SentAt.Minute())

	return fmt.Sprintf("[%s] [%s] (from: %s, subject: %s, reason: %s): %s",
		sentAtStr, m.In.Name, m.Sender.Name, m.Subject.Name,
		m.Reason, m.Text)
}

// New creates a new Msg
func NewMsg(in *Source, sender *Source, sentAt *time.Time, subject *Source,
	reason string, text string) *Msg {

	return &Msg{
		In:      in,
		Sender:  sender,
		SentAt:  sentAt,
		Subject: subject,
		Reason:  reason,
		Text:    text,
	}
}
