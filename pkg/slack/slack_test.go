package slack

import (
	"fmt"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestPost(t *testing.T) {
	log := logrus.StandardLogger()
	log.SetLevel(logrus.DebugLevel)
	url := "https://hooks.slack.com/services/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"
	payload := &slack.WebhookMessage{
		Text: "hello world",
		Attachments: []slack.Attachment{
			{
				AuthorID:   "matt",
				Text:       fmt.Sprintf(":) :p0_912: [dont-click-me](http://www.google.com.tw)"),
				MarkdownIn: []string{"dont-click-me"},
			},
		},
	}
	fmt.Sprintf("%+v\n", url)
	fmt.Sprintf("%+v\n", payload)
	//resp := slack.PostWebhook(url, payload)
	//fmt.Printf("%+v", resp)
}
