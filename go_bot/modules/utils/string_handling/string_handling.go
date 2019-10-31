/*
 *   Copyright 2019 ATechnoHazard  <amolele@gmail.com>
 *
 *   Permission is hereby granted, free of charge, to any person obtaining a copy
 *   of this software and associated documentation files (the "Software"), to deal
 *   in the Software without restriction, including without limitation the rights
 *   to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *   copies of the Software, and to permit persons to whom the Software is
 *   furnished to do so, subject to the following conditions:
 *
 *   The above copyright notice and this permission notice shall be included in all
 *   copies or substantial portions of the Software.
 *
 *   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *   IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *   FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *   AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *   LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *   OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *   SOFTWARE.
 */

package string_handling

import (
	"fmt"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/error_handling"
	"github.com/PaulSonOfLars/gotgbot/ext"
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
			error_handling.HandleErr(err)
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
		error_handling.HandleErr(err)
		return -1
	}
}

func FormatText(format string, args ...string) string {
	r := strings.NewReplacer(args...)
	return r.Replace(format)
}
