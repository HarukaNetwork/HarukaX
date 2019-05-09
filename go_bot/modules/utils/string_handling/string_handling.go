package string_handling

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/atechnohazard/ginko/go_bot/modules/utils/error_handling"
	"strconv"
	"strings"
	"time"
)

func ExtractTime(m *ext.Message, timeVal string) int64 {
	lastLetter := timeVal[len(timeVal)-1:]
	lastLetter = strings.ToLower(lastLetter)
	var banTime int64
	if strings.ContainsAny(lastLetter, "m & d & h") {
		t := timeVal[:len(timeVal)-1]
		timeNum, err := strconv.Atoi(t)
		if err != nil {
			_, err := m.ReplyText("Invalid time amount specified.")
			error_handling.HandleErrorGracefully(err)
			return -1
		}

		if lastLetter == "m" {
			banTime = time.Now().Unix() + int64(timeNum*60)
		} else if lastLetter == "h" {
			banTime = time.Now().Unix() + int64(timeNum*60*60)
		} else if lastLetter == "d" {
			banTime = time.Now().Unix() + int64(timeNum*24*60*60)
		} else {
			return 0
		}
		return banTime
	} else {
		_, err := m.ReplyText(fmt.Sprintf("Invalid time type specified. Expected m, h, or d got: %s", lastLetter))
		error_handling.HandleErrorGracefully(err)
		return -1
	}
}

func FormatText(format string, args ...string) string {
	r := strings.NewReplacer(args...)
	return r.Replace(format)
}
