package main

import (
	"context"
	"log"
	"os"

	"github.com/argally/helpbot/internal/slackbot"

	"github.com/joho/godotenv"
	"github.com/slack-io/slacker"
)

func main() {
	godotenv.Load(".env")

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	bot.AddCommand(slackbot.AzureHelpers())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
