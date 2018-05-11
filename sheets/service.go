package sheets

import "github.com/Noah-Huppert/should-have-worked-from-home/config"
import "gopkg.in/Iwark/spreadsheet.v2"
import "golang.org/x/oauth2/google"
import "io/ioutil"
import "context"
import "fmt"

// NewService creates a new Google API Spreadsheet service.
func NewService(ctx context.Context, cfg *config.Config) (*spreadsheet.Service,
	error) {

	// Load key file
	keyDat, err := ioutil.ReadFile(cfg.GAPIKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error loading GAPI key file: %s",
			err.Error())
	}

	// Get GAPI OAuth configuration used for retrieving API token
	jwtConf, err := google.JWTConfigFromJSON(keyDat, spreadsheet.Scope)
	if err != nil {
		return nil, fmt.Errorf("error loading GAPI OAuth "+
			"configuration: %s", err.Error())
	}

	// Create client
	client := jwtConf.Client(ctx)

	// Create spreadsheet API service
	svc := spreadsheet.NewServiceWithClient(client)

	return svc, nil
}
