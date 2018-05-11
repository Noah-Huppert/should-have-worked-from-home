package gapi

import "github.com/Noah-Huppert/should-have-worked-from-home/config"
import "golang.org/x/oauth2/google"
import "golang.org/x/oauth2"
import "io/ioutil"
import "fmt"

// GAPISpreadsheetScope is the GAPI OAuth spreadsheet scope
const GAPISpreadsheetScope string = "https://www.googleapis.com/auth/spreadsheets"

// GetTokenURL retrieves a Google API authentication URL for retrieving a user
// OAuth token.
func GetTokenURL(cfg *config.Config) (string, error) {
	// Read client secret file
	bytes, err := ioutil.ReadFile(cfg.GAPIPrivateKeyPath)
	if err != nil {
		return "", fmt.Errorf("error reading GAPI private key file: "+
			"%s", err.Error())
	}

	// Make OAuth token config
	tokenConfig, err := google.ConfigFromJSON(bytes, GAPISpreadsheetScope)
	if err != nil {
		return "", fmt.Errorf("error creating OAuth token "+
			"configuration: %s", err.Error())
	}

	// Get url
	authURL := tokenConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return authURL, nil
}
