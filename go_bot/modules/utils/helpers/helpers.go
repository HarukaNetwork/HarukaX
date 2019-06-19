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
				smallMsg += line + "\n"
			} else {
				result = append(result, smallMsg)
				smallMsg = line + "\n"
			}
		}
		result = append(result, smallMsg)
		return result
	}
}

func SplitQuotes(text string) []string {
	smartOpen := "“"
	smartClose := "”"
	startChars := []string{"'", `"`, smartOpen}
	broke := false

	for _, char := range startChars {
		if strings.HasPrefix(text, char) {
			counter := 1
			for counter < len(text) {
				if text[counter] == '\\' {
					counter++
				} else if text[counter] == text[0] || (string(text[0]) == smartOpen && string(text[counter]) == smartClose) {
						broke = true
						break
				}
				counter++
			}
			if !broke {
				return strings.SplitN(text, " ", 2)
			}

			key := RemoveEscapes(strings.TrimSpace(text[1:counter]))
			rest := strings.TrimSpace(text[counter + 1:])

			if key == "" {
				key = string(text[0]) + string(text[0])
			}
			tmp := make([]string, 2)
			tmp[0] = key
			tmp[1] = rest
			return tmp
		}
	}
	return strings.SplitN(text, " ", 2)
}

func RemoveEscapes(text string) string {
	counter := 0
	res := ""
	isEscaped := false

	for counter < len(text) {
		if isEscaped {
			res += string(text[counter])
			isEscaped = false
		} else if text[counter] == '\\' {
			isEscaped = true
		} else {
			res += string(text[counter])
		}
		counter++
	}
	return res
}