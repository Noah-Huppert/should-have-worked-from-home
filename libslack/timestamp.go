package libslack

import "strings"
import "strconv"
import "time"
import "fmt"

// ConvertTimestamp converts a Slack style timestamp to a Golang time.
// Slack timestamps are in the format:
//
// 	xxxxxxxxxx.yyyyy
//
// The portion before the decimal is a unix timestamp. The portion after the
// decimal is a number used to ensure correct ordering of Slack message.
func ConvertTimestamp(ts string) (*time.Time, error) {
	// Split
	parts := strings.Split(ts, ".")
	unixTsStr := parts[0]

	// Part into number
	unixTs, err := strconv.ParseInt(unixTsStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error converting unix timestamp "+
			"portion into integer, unix timestamp portion: %s, "+
			"error: %s", unixTsStr, err.Error())
	}

	// Parse into time
	t := time.Unix(unixTs, 0)
	return &t, nil
}
