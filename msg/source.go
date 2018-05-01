package msg

import "fmt"

// Source indicates where a Slack message came from
type Source struct {
	// IsUser indicates the source is a user
	IsUser bool

	// IsChannel indicates the source is a channel
	IsChannel bool

	// ID is the Slack ID of the source
	ID string

	// Name is the human readable name of the source
	Name string
}

// String returns a text representation of a Source
func (s Source) String() string {
	prefix := "user"
	if s.IsChannel {
		prefix = "channel"
	}

	return fmt.Sprintf("%s:%s (id: %s)", prefix, s.Name, s.ID)
}

// NewUserSource creates a new Source for a user
func NewUserSource(id string, name string) *Source {
	return &Source{
		IsUser:    true,
		IsChannel: false,
		ID:        id,
		Name:      name,
	}
}

// NewChannelSource creates a new Source for a channel
func NewChannelSource(id string, name string) *Source {
	return &Source{
		IsUser:    false,
		IsChannel: true,
		ID:        id,
		Name:      name,
	}
}
