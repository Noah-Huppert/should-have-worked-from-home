package gapi

import "github.com/Noah-Huppert/should-have-worked-from-home/config"
import "github.com/Noah-Huppert/should-have-worked-from-home/stdin"
import "golang.org/x/oauth2/google"
import "golang.org/x/oauth2"
import "io/ioutil"
import "fmt"
import "log"
import "context"
import "encoding/json"
import "strings"

// GAPISpreadsheetScope is the GAPI OAuth spreadsheet scope
const GAPISpreadsheetScope string = "https://www.googleapis.com/auth/spreadsheets"

// GetOAuthConfig constructs a OAuth configuration from a GAPI private key
// file.
func GetOAuthConfig(cfg *config.Config) (*oauth2.Config, error) {
	// Read client secret file
	bytes, err := ioutil.ReadFile(cfg.GAPIPrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error reading GAPI private key file: "+
			"%s", err.Error())
	}

	// Make OAuth config
	gapiOAuthConfig, err := google.ConfigFromJSON(bytes, GAPISpreadsheetScope)
	if err != nil {
		return nil, fmt.Errorf("error creating OAuth configuration:"+
			" %s", err.Error())
	}

	return gapiOAuthConfig, nil
}

// GetTokenURL retrieves a Google API authentication URL for retrieving a user
// OAuth token.
func GetTokenURL(cfg *config.Config, gapiOAuthConfig *oauth2.Config) (string, error) {
	// Get url
	authURL := gapiOAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return authURL, nil
}

// ReadAuthorizationCode presents the user with instructions to retrieve an
// GAPI OAuth code. Then this authorization token is exchanged for an access
// token.
func ReadAuthorizationCode(ctx context.Context, cfg *config.Config,
	logger *log.Logger, gapiOAuthConfig *oauth2.Config) error {

	// Get authorization token URL
	authURL, err := GetTokenURL(cfg, gapiOAuthConfig)
	if err != nil {
		return fmt.Errorf("error retrieving GAPI authorization "+
			"token URL: %s", err.Error())
	}

	// Prompt for authorization token
	logger.Printf("No GAPI access token provided. \nPlease navigate "+
		"to the following URL and enter the code provided in the "+
		"console: \n\n%s\n\n", authURL)

	authCode, err := stdin.Prompt(logger, "GAPI authorization code: ")
	if err != nil {
		return fmt.Errorf("error reading GAPI authorization code from "+
			"terminal: %s", err.Error())
	}

	// Exchange authorization code for access token
	accessToken, err := gapiOAuthConfig.Exchange(ctx, authCode)
	if err != nil {
		return fmt.Errorf("error exchanging authorization code for "+
			"access token: %s", err.Error())
	}

	// Convert access token to JSON text
	accessTokenTxt, err := json.Marshal(accessToken)
	if err != nil {
		return fmt.Errorf("error converting GAPI access token into "+
			"JSON: %s", err.Error())
	}

	// Instruct user to save GAPI access token
	logger.Println()
	logger.Printf("GAPI access token retrieved, please save the "+
		"following value in the \"%s\" environment variable: \n\n%s",
		config.EnvKeyGAPIAccessToken, accessTokenTxt)

	return nil
}

// MakeToken creates an oauth2.Token object from a string
func MakeToken(tokenStr string) (*oauth2.Token, error) {
	gapiAccessToken := &oauth2.Token{}

	// Create io.Reader from token string
	tokenReader := strings.NewReader(tokenStr)

	// Decode token into oauth2.Token var
	decoder := json.NewDecoder(tokenReader)

	err := decoder.Decode(gapiAccessToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode access token string "+
			"into oauth2.Token: %s", err.Error())
	}

	return gapiAccessToken, nil
}
