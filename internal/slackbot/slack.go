package slackbot

import "github.com/slack-io/slacker"

func FetchUserDetails(botCtx *slacker.CommandContext) string {
	userID := botCtx.Event().UserID
	user, err := botCtx.SlackClient().GetUserInfo(userID)
	if err != nil {
		return "Unknown"
	}
	return user.RealName

}
