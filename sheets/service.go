package sheets

import "github.com/Noah-Huppert/should-have-worked-from-home/config"
import "github.com/Noah-Huppert/should-have-worked-from-home/gapi"
import "gopkg.in/Iwark/spreadsheet.v2"
import "golang.org/x/oauth2"
import "context"
import "fmt"

// NewService creates a new Google API Spreadsheet service.
func NewService(ctx context.Context, cfg *config.Config,
	gapiOAuthConfig *oauth2.Config) (*spreadsheet.Service, error) {

	// Make OAuth token
	token, err := gapi.MakeToken(cfg.GAPIAccessToken)
	if err != nil {
		return nil, fmt.Errorf("error creating GAPI access token: %s",
			err.Error())
	}

	// Make GAPI client
	client := gapiOAuthConfig.Client(ctx, token)

	// Create spreadsheet API service
	svc := spreadsheet.NewServiceWithClient(client)

	return svc, nil
}
