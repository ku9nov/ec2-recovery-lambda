package actions

import (
	"os"
	"strconv"
	"time"

	"github.com/slack-go/slack"
)

const (
	AlarmMessage   string = "Instance is broken."
	RestartMessage string = "Succesfully restarted instance."
	GreenColor     string = "#36a64f"
	RedColor       string = "#FF0000"
)

func SendMessageToSlack(text, color, instanceName, instanceID string, instanceCount int) {
	var instanceValue string
	if os.Getenv("SLACK_ENABLE") == "true" {
		token := os.Getenv("SLACK_AUTH_TOKEN")
		channelID := os.Getenv("SLACK_CHANNEL_ID")
		client := slack.New(token, slack.OptionDebug(true))
		switch instanceName {
		case "":
			instanceValue = instanceID
		default:
			instanceValue = instanceName
		}
		attachment := slack.Attachment{
			Title: "EC2 recovery Lambda Notification",
			Text:  text,
			Color: color,
			Fields: []slack.AttachmentField{
				{
					Title: "Host",
					Value: instanceValue,
					Short: true,
				},
				{
					Title: "Total number of instances",
					Value: strconv.Itoa(instanceCount),
					Short: true,
				},
				{
					Title: "Event time",
					Value: time.Now().Format("2006.01.02 15:04:05"),
					Short: true,
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
