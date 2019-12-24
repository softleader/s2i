package slack

import "github.com/nlopes/slack"

// Post post webhook to url with text
func Post(url, text string) error {
	payload := &slack.WebhookMessage{
		Text: text,
	}
	return slack.PostWebhook(url, payload)
}
