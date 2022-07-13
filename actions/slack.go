package actions

import (
	"os"
	"time"

	"github.com/slack-go/slack"
)

const (
	AlarmMessage   string = "Instance is broken."
	RestartMessage string = "Succesfully restarted instance."
	GreenColor     string = "#36a64f"
	RedColor       string = "#FF0000"
)

func SendMessageToSlack(text, color, instanceID string) {
	if os.Getenv("SLACK_ENABLE") == "true" {
		token := os.Getenv("SLACK_AUTH_TOKEN")
		channelID := os.Getenv("SLACK_CHANNEL_ID")
		client := slack.New(token, slack.OptionDebug(true))
		attachment := slack.Attachment{
			Title: "EC2 recovery Lambda Notification",
			Text:  text,
			Color: color,
			Fields: []slack.AttachmentField{
				{
					Title: "Host",
					Value: instanceID,
				},
				{
					Title: "Event time",
					Value: time.Now().Format("2006.01.02 15:04:05"),
				},
			},
		}
		_, timestamp, err := client.PostMessage(
			channelID,
			slack.MsgOptionAttachments(attachment),
		)
		_ = timestamp
		if err != nil {
			panic(err)
		}
	}

}
