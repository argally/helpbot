package slackbot

import (
	"fmt"

	"github.com/argally/helpbot/internal/azure"
	"github.com/slack-go/slack"
	"github.com/slack-io/slacker"
)

// AzureHelpers defines a Slack command for creating Azure resources.
//
// Command:
//
//	azure create {resource-type} {resource-name} {location}
//
// Description:
//
//	This command allows users to create new Azure resources by specifying the
//	resource type, resource name, and location. If the location is not provided,
//	it defaults to "eastus".
//
// Examples:
//
//	Supported resources:
//	  - azure create blob-container mysampleblob eastus
//	  - azure create storage-account mysamplestorage westus
//
// Parameters:
//   - resource-type: The type of Azure resource to create (e.g., blob-container, storage-account).
//   - resource-name: The name of the Azure resource to create.
//   - location: (Optional) The Azure region where the resource will be created. Defaults to "eastus".
//
// Behavior:
//   - Validates the input parameters and ensures required fields are provided.
//   - Calls the AzureServicesCreate function to create the resource.
//   - Responds with an error message if the creation fails.
//   - Sends a Slack message with the result of the resource creation if successful.
func AzureHelpers() *slacker.CommandDefinition {
	azurebot := &slacker.CommandDefinition{
		Command:     "azure create {resource-type} {resource-name} {location}",
		Description: "Create new Azure Resource",
		Examples: []string{
			"Supported resources:",
			"azure create blob-container mysampleblob eastus",
			"azure create storage-account mysamplestorage westus",
		},

		Handler: func(ctx *slacker.CommandContext) {
			//resourcetype := ctx.Request().Param("resource-type")

			resourceType := ctx.Request().Param("resource-type")
			resourceName := ctx.Request().Param("resource-name")
			location := ctx.Request().Param("location")

			// Validate input parameters
			if resourceType == "" || resourceName == "" {
				ctx.Response().Reply(":x: Missing required parameters. Please provide `resource-type`, `resource-name`")
				return
			}
			if location == "" {
				ctx.Response().Reply(":information_source: Location will default to eastus")
			}

			result, err := azure.AzureServicesCreate(ctx.Context(), resourceName, resourceType, location)
			if err != nil {
				ctx.Response().Reply(fmt.Sprintf(":x: Error creating resource %v", err))
				return
			}
			attachments := createSlackAttachment(result)
			ctx.Response().PostBlocks(ctx.Event().ChannelID, attachments)
		},
	}
	return azurebot
}

func createSlackAttachment(message string) []slack.Block {
	return []slack.Block{
		slack.NewDividerBlock(),
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", message, false, false),
			nil,
			nil,
		),
		slack.NewDividerBlock(),
	}
}
