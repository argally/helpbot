package main

import (
	"context"
	"log"
	"os"

	"github.com/argally/helpbot/internal/slackbot"
	"github.com/slack-go/slack"

	"github.com/joho/godotenv"
	"github.com/slack-io/slacker"
)

func main() {
	godotenv.Load(".env") //nolint:errcheck

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	bot.AddCommand(slackbot.AzureHelpers())
	bot.AddCommand(&slacker.CommandDefinition{
		Command: "hello",
		Handler: func(ctx *slacker.CommandContext) {
			// Get the native Slack client
			slackClient := ctx.SlackClient()

			// Iterate over the channels and post a message to each
			for _, channel := range slackbot.ChannelList {
				_, _, err := slackClient.PostMessage(channel.ID, slack.MsgOptionText(ctx.Event().ChannelID, false))
				if err != nil {
					log.Printf("Failed to post message to channel %s: %v", channel, err)
				} else {
					log.Printf("Message successfully sent to channel %s", channel)
				}
			}
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
