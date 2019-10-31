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

package users

import (
	"github.com/ATechnoHazard/ginko/go_bot/modules/sql"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/handlers/Filters"
	"log"
	"strconv"
)

func logUsers(_ ext.Bot, u *gotgbot.Update) error {
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	sql.UpdateUser(msg.From.Id,
		msg.From.Username,
		strconv.Itoa(chat.Id),
		chat.Title)

	if msg.ReplyToMessage != nil {
		sql.UpdateUser(msg.From.Id,
			msg.From.Username,
			strconv.Itoa(chat.Id),
			chat.Title)
	}

	if msg.ForwardFrom != nil {
		sql.UpdateUser(msg.ForwardFrom.Id,
			msg.ForwardFrom.Username, "nil", "nil")
	}

	return gotgbot.ContinueGroups{}
}

func GetUserId(username string) int {
	if len(username) <= 5 {
		return 0
	}
	if username[0] == '@' {
		username = username[1:]
	}
	users := sql.GetUserIdByName(username)
	if users == nil {
		return 0
	}

	return users.UserId
}

func LoadUsers(u *gotgbot.Updater) {
	defer log.Println("Loading module users")
	u.Dispatcher.AddHandler(handlers.NewMessage(Filters.All, logUsers))
}
