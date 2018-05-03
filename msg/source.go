package msg

import "fmt"

// SourceType represents different types of Slack resources
type SourceType string

const (
	User    SourceType = "user"
	Channel SourceType = "channel"
	IM      SourceType = "im"
)

// SourcesBotSelfKey is the special key the bot's own identity will be stored
// under in the sources cache.
//
// This allows one to check the sources cache for the bot's identity without
// knowing the id of the bot (As we can only know the id of the bot by calling
// libslack.GetIdentity).
const SourcesBotSelfKey string = "self"

// Source indicates where a Slack message came from
type Source struct {
	// ID is the Slack ID of the source
	ID string

	// Name is the human readable name of the source
	Name string

	// Type is the type of Slack resource the source indicates
	Type SourceType
}

// String returns a text representation of a Source
func (s Source) String() string {
	return fmt.Sprintf("%s: %s (id: %s)", s.Type, s.Name, s.ID)
}

// NewSource creates a new Source instance
func NewSource(id string, name string, t SourceType) *Source {
	return &Source{
		ID:   id,
		Name: name,
		Type: t,
	}
}
