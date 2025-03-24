package slackbot

type Channel struct {
	ID   string
	Name string
}

var Random = Channel{ID: "CQVRAQSNM", Name: "random"}
var Automation = Channel{ID: "CR3PKAYJ1", Name: "automation"}

var ChannelList = []Channel{
	Random,
	Automation,
}
