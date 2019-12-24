package slack

import "testing"

func TestPost(t *testing.T) {
	Post("https://hooks.slack.com/services/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", "hello world")
}
