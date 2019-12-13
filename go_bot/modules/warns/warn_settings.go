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

package warns

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/ATechnoHazard/ginko/go_bot/modules/sql"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/chat_status"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/error_handling"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
)

func setWarnLimit(_ ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat
	msg := u.EffectiveMessage
	user := u.EffectiveUser

	// Check permissions
	if !chat_status.RequireUserAdmin(chat, msg, user.Id) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.RequireBotAdmin(chat, msg) {
		return gotgbot.EndGroups{}
	}

	if len(args) > 0 {
		for _, elem := range args[0] {
			if unicode.IsDigit(elem) {
				num, err := strconv.Atoi(args[0])
				error_handling.HandleErr(err)

				if num < 3 {
					_, err := msg.ReplyText("The minimum warn limit is 3!")
					return err
				} else {
					go sql.SetWarnLimit(strconv.Itoa(chat.Id), num)
					_, err := msg.ReplyHTML(fmt.Sprintf("Updated the warn limit to <b>%v</b>", num))
					return err
				}
			}
		}
		_, err := msg.ReplyText("Give me a number as an argument!")
		return err
	} else {
		limit, softWarn := sql.GetWarnSetting(strconv.Itoa(chat.Id))
		_, err := msg.ReplyHTML(fmt.Sprintf("The current warn limit is <b>%v</b>.\nThe soft warn setting is set to <b>%v</b>.", limit, softWarn))
		return err
	}
}

func setWarnStrength(_ ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat
	msg := u.EffectiveMessage
	user := u.EffectiveUser

	// Check permissions
	if !chat_status.RequireUserAdmin(chat, msg, user.Id) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.RequireBotAdmin(chat, msg) {
		return gotgbot.EndGroups{}
	}

	if len(args) > 0 {
		if strings.ToLower(args[0]) == "on" || strings.ToLower(args[0]) == "yes" {
			go sql.SetWarnStrength(strconv.Itoa(chat.Id), false)
			_, err := msg.ReplyText("Too many warns will now result in a ban!")
			return err
		} else if strings.ToLower(args[0]) == "off" || strings.ToLower(args[0]) == "no" {
			go sql.SetWarnStrength(strconv.Itoa(chat.Id), true)
			_, err := msg.ReplyText("Too many warns will now result in a kick! User will be able to join again after.")
			return err
		} else {
			_, err := msg.ReplyText("I only understand on/yes/no/off!")
			return err
		}
	} else {
		limit, softWarn := sql.GetWarnSetting(strconv.Itoa(chat.Id))
		if softWarn {
			_, err := msg.ReplyHTML(fmt.Sprintf("Warns are currently set to <b>kick</b> users when they exceed <b>%v</b> warns.", limit))
			return err
		} else {
			_, err := msg.ReplyHTML(fmt.Sprintf("Warns are currently set to <b>ban</b> users when they exceed <b>%v</b> warns.", limit))
			return err
		}
	}
}
