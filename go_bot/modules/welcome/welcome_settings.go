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

package welcome

import (
	"strconv"
	"strings"

	"github.com/ATechnoHazard/ginko/go_bot/modules/sql"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/chat_status"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/helpers"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
)

func welcome(bot ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat

	if !chat_status.IsUserAdmin(chat, u.EffectiveUser.Id) {
		_, _ = u.EffectiveMessage.ReplyText("You need to be an admin to do this.")
		return gotgbot.ContinueGroups{}
	}

	if len(args) == 0 || strings.ToLower(args[0]) == "noformat" {
		noformat := len(args) > 0 && strings.ToLower(args[0]) == "noformat"
		welcPrefs := sql.GetWelcomePrefs(strconv.Itoa(chat.Id))
		_, _ = u.EffectiveMessage.ReplyHTMLf("I am currently welcoming users: <code>%v</code>"+
			"\nI am currently deleting old welcomes: <code>%v</code>"+
			"\nI am currently deleting service messages: <code>%v</code>"+
			"\nOn joining, I am currently muting users: <code>%v</code>"+
			"\nThe welcome message not filling the {} is:",
			welcPrefs.ShouldWelcome,
			welcPrefs.CleanWelcome != 0,
			welcPrefs.DelJoined,
			welcPrefs.ShouldMute)

		if welcPrefs.WelcomeType == sql.BUTTON_TEXT {
			buttons := sql.GetWelcomeButtons(strconv.Itoa(chat.Id))
			if noformat {
				welcPrefs.CustomWelcome += helpers.RevertButtons(buttons)
				_, err := u.EffectiveMessage.ReplyHTML(welcPrefs.CustomWelcome)
				return err
			} else {
				keyb := helpers.BuildWelcomeKeyboard(buttons)
				keyboard := ext.InlineKeyboardMarkup{InlineKeyboard: &keyb}
				send(bot, u, welcPrefs.CustomWelcome, &keyboard, sql.DefaultWelcome)
			}
		} else {
			_, err := EnumFuncMap[welcPrefs.WelcomeType](bot, chat.Id, welcPrefs.CustomWelcome) // needs change
			return err
		}
	} else if len(args) >= 1 {
		switch strings.ToLower(args[0]) {
		case "on", "yes":
			go sql.SetWelcPref(strconv.Itoa(chat.Id), true)
			_, err := u.EffectiveMessage.ReplyText("I'll welcome users from now on.")
			return err
		case "off", "no":
			go sql.SetWelcPref(strconv.Itoa(chat.Id), false)
			_, err := u.EffectiveMessage.ReplyText("I'll not welcome users from now on.")
			return err
		default:
			_, err := u.EffectiveMessage.ReplyText("I understand 'on/yes' or 'off/no' only!")
			return err
		}
	}
	return nil
}

func setWelcome(_ ext.Bot, u *gotgbot.Update) error {
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	if !chat_status.IsUserAdmin(chat, u.EffectiveUser.Id) {
		_, _ = u.EffectiveMessage.ReplyText("You need to be an admin to do this.")
		return gotgbot.ContinueGroups{}
	}

	text, dataType, content, buttons := helpers.GetWelcomeType(msg)
	if dataType == -1 {
		_, err := msg.ReplyText("You didn't specify what to reply with!")
		return err
	}

	btns := make([]sql.WelcomeButton, len(buttons))
	for i, btn := range buttons {
		btns[i] = sql.WelcomeButton{
			ChatId:   strconv.Itoa(chat.Id),
			Name:     btn.Name,
			Url:      btn.Content,
			SameLine: btn.SameLine,
		}
	}

	if text != "" {
		go sql.SetCustomWelcome(strconv.Itoa(chat.Id), text, btns, dataType)
	} else {
		go sql.SetCustomWelcome(strconv.Itoa(chat.Id), content, btns, dataType)
	}

	_, err := msg.ReplyText("Successfully set custom welcome message!")
	return err
}

func resetWelcome(_ ext.Bot, u *gotgbot.Update) error {
	chat := u.EffectiveChat
	if !chat_status.IsUserAdmin(chat, u.EffectiveUser.Id) {
		_, _ = u.EffectiveMessage.ReplyText("You need to be an admin to do this.")
		return gotgbot.ContinueGroups{}
	}

	go sql.SetCustomWelcome(strconv.Itoa(chat.Id), sql.DefaultWelcome, nil, sql.TEXT)

	_, err := u.EffectiveMessage.ReplyText("Succesfully reset custom welcome message to default!")
	return err
}

func cleanWelcome(bot ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat
	if !chat_status.IsUserAdmin(chat, u.EffectiveUser.Id) {
		_, _ = u.EffectiveMessage.ReplyText("You need to be an admin to do this.")
		return gotgbot.ContinueGroups{}
	}

	if len(args) == 0 {
		cleanPref := sql.GetCleanWelcome(strconv.Itoa(chat.Id))
		if cleanPref != 0 {
			_, err := u.EffectiveMessage.ReplyText("I should be deleting welcome messages up to two days old.")
			return err
		} else {
			_, err := u.EffectiveMessage.ReplyText("I'm currently not deleting old welcome messages!")
			return err
		}
	} else {
		switch strings.ToLower(args[0]) {
		case "off", "no":
			go sql.SetCleanWelcome(strconv.Itoa(chat.Id), 0)
			_, err := u.EffectiveMessage.ReplyText("I'll try to delete old welcome messages!")
			return err
		case "on", "yes":
			go sql.SetCleanWelcome(strconv.Itoa(chat.Id), 1)
			_, err := u.EffectiveMessage.ReplyText("I'll try to delete old welcome messages!")
			return err
		default:
			_, err := u.EffectiveMessage.ReplyText("I understand 'on/yes' or 'off/no' only!")
			return err
		}
	}
}
