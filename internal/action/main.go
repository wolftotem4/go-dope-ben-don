package action

import (
	"errors"
	"html"
	"regexp"
	"strings"
)

var ErrLoginRequired = errors.New("login required")
var ErrPageExpired = errors.New("page expired")

func handleFeedbackErrMsg(content []byte) error {
	regex := regexp.MustCompile(`(?si)<span class="feedbackPanelERROR">(.*?)</span>`)
	match := regex.FindSubmatch(content)
	if len(match) >= 2 {
		return errors.New(html.UnescapeString(strings.TrimSpace(string(match[1]))))
	}
	return nil
}
