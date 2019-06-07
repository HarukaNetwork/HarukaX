package helpers

import (
	"fmt"
	"html"
	"strings"
)

var MAX_MESSAGE_LENGTH = 4096

func MentionHtml(userId int, name string) string {
	return fmt.Sprintf("<a href=\"tg://user?id=%d\">%s</a>", userId, html.EscapeString(name))
}

func SplitMessage(msg string) []string {
	if len(msg) > MAX_MESSAGE_LENGTH {
		tmp := make([]string, 1)
		tmp[0] = msg
		return tmp
	} else {
		lines := strings.Split(msg, "\n")
		smallMsg := ""
		result := make([]string, 0)
		for _, line := range lines {
			if len(smallMsg) + len(line) < MAX_MESSAGE_LENGTH {
				smallMsg += line
			} else {
				result = append(result, smallMsg)
				smallMsg = line
			}
		}
		result = append(result, smallMsg)
		return result
	}
}