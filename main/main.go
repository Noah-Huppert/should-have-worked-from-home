package main

import "github.com/Noah-Huppert/should-have-worked-from-home/config"
import "github.com/Noah-Huppert/should-have-worked-from-home/bot"
import "github.com/Noah-Huppert/should-have-worked-from-home/sheets"
import "github.com/Noah-Huppert/should-have-worked-from-home/gapi"
import "os"
import "golang.org/x/oauth2"
import "os/signal"
import "syscall"
import "log"
import "context"
import "fmt"

func main() {
	logger := log.New(os.Stdout, "server: ", 0)

	// Signal handler
	signals := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())

	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		for _ = range signals {
			fmt.Println("")
			logger.Println("received termination signal")
			cancel()
		}
	}()

	// Load config
	cfg, err := config.New()
	var gapiOAuthConfig *oauth2.Config

	// Get authorization token from user
	if err == config.ErrNoGAPIAccessToken {
		// Get GAPI OAuth config
		gapiOAuthConfig, err := gapi.GetOAuthConfig(cfg)
		if err != nil {
			logger.Fatalf("error retrieving GAPI OAuth "+
				"configuration: %s", err.Error())
			return
		}

		// Get GAPI access token
		err = gapi.ReadAuthorizationCode(ctx, cfg, logger, gapiOAuthConfig)
		if err != nil {
			logger.Fatalf("error reading GAPI authorization "+
				"code: %s", err.Error())
		}

		// User should relaunch program with GAPI access token
		// environment variable
		return
	}
	if err != nil {
		logger.Fatalf("error loading configuration: %s", err.Error())
	}

	// Make Sheet client
	svc, err := sheets.NewService(ctx, cfg, gapiOAuthConfig)
	if err != nil {
		logger.Fatalf("error creating GAPI Sheets service: %s",
			err.Error())
	}

	sheet, err := sheets.NewSheet(svc, cfg.SpreadsheetID,
		cfg.SpreadsheetPageName)
	if err != nil {
		logger.Fatalf("error creating sheet: %s", err.Error())
	}
	logger.Println(sheet)

	// Start listener
	msgs, errs := bot.Listen(ctx, cfg)

	for {
		select {
		case err := <-errs:
			// If exit signal
			if err == nil {
				logger.Fatalln("listener exiting")
			} else {
				logger.Printf("listener error: %s\n", err.Error())
			}

		case msg := <-msgs:
			logger.Printf("received message: %s", msg)
			// TODO: Use Google Sheets API to record result
			// TODO: Send response back to subject
		}
	}
}
