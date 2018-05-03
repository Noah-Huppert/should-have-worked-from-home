package main

import "github.com/Noah-Huppert/should-have-worked-from-home/config"
import "github.com/Noah-Huppert/should-have-worked-from-home/bot"
import "os"
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
	if err != nil {
		logger.Fatalf("error loading configuration: %s", err.Error())
	}

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
