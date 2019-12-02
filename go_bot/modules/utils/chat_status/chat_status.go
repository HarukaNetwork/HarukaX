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

package chat_status

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ATechnoHazard/ginko/go_bot"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/caching"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/error_handling"
	"github.com/PaulSonOfLars/gotgbot/ext"
)

type Cache struct {
	Admin []string `json:"admin"`
}

func CanDelete(chat *ext.Chat, botId int) bool {
	k, err := chat.GetMember(botId)
	error_handling.HandleErr(err)
	return k.CanDeleteMessages
}

func IsUserBanProtected(chat *ext.Chat, userId int, member *ext.ChatMember) bool {
	if chat.Type == "private" || contains(go_bot.BotConfig.SudoUsers, strconv.Itoa(userId)) {
		return true
	}
	if member == nil {
		mem, err := chat.GetMember(userId)
		error_handling.HandleErr(err)
		member = mem
	}
	if member.Status == "administrator" || member.Status == "creator" {
		return true
	} else {
		return false
	}
}

func IsUserAdmin(chat *ext.Chat, userId int) bool {
	if chat.Type == "private" || contains(go_bot.BotConfig.SudoUsers, strconv.Itoa(userId)) {
		return true
	}

	admins, _ := caching.CACHE.Get(fmt.Sprintf("admin_%v", chat.Id))
	go cacheAdmins(chat)

	var adminList Cache
	_ = json.Unmarshal(admins, &adminList)
	return contains(adminList.Admin, strconv.Itoa(userId))
}

func IsBotAdmin(chat *ext.Chat, member *ext.ChatMember) bool {
	if chat.Type == "private" {
		return true
	}
	if member == nil {
		mem, err := chat.GetMember(chat.Bot.Id)
		error_handling.HandleErr(err)
		if mem == nil {
			return false
		}
		member = mem

	}
	if member.Status == "administrator" || member.Status == "creator" {
		return true
	} else {
		return false
	}
}

func RequireBotAdmin(chat *ext.Chat, msg *ext.Message) bool {
	if !IsBotAdmin(chat, nil) {
		_, err := msg.ReplyText("I'm not admin!")
		error_handling.HandleErr(err)
		return false
	}
	return true
}

func RequireUserAdmin(chat *ext.Chat, msg *ext.Message, userId int, member *ext.ChatMember) bool {
	if !IsUserAdmin(chat, userId) {
		_, err := msg.ReplyText("You must be an admin to perform this action.")
		error_handling.HandleErr(err)
		return false
	}
	return true
}

func IsUserInChat(chat *ext.Chat, userId int) bool {
	member, err := chat.GetMember(userId)
	error_handling.HandleErr(err)
	if member.Status == "left" || member.Status == "kicked" {
		return false
	} else {
		return true
	}
}

func CanPromote(bot ext.Bot, chat *ext.Chat) bool {
	botChatMember, err := chat.GetMember(bot.Id)
	error_handling.HandleErr(err)
	if !botChatMember.CanPromoteMembers {
		_, err := bot.SendMessage(chat.Id, "I can't promote/demote people here! Make sure I'm admin and can appoint new admins.")
		error_handling.HandleErr(err)
		return false
	}
	return true
}

func CanPin(bot ext.Bot, chat *ext.Chat) bool {
	botChatMember, err := chat.GetMember(bot.Id)
	error_handling.HandleErr(err)
	if !botChatMember.CanPinMessages {
		_, err := bot.SendMessage(chat.Id, "I can't pin messages here! Make sure I'm admin and can pin messages.")
		error_handling.HandleErr(err)
		return false
	}
	return true
}

func CanRestrict(bot ext.Bot, chat *ext.Chat) bool {
	botChatMember, err := chat.GetMember(bot.Id)
	error_handling.HandleErr(err)
	if !botChatMember.CanRestrictMembers {
		_, err := bot.SendMessage(chat.Id, "I can't restrict people here! Make sure I'm admin and can appoint new admins.")
		error_handling.HandleErr(err)
		return false
	}
	return true
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func cacheAdmins(chat *ext.Chat) {
	adminList, err := chat.GetAdministrators()
	error_handling.HandleErr(err)
	admins := make([]string, 0)

	for _, admin := range adminList {
		admins = append(admins, strconv.Itoa(admin.User.Id))
	}

	adminCache := &Cache{admins}
	adminJson, _ := json.Marshal(adminCache)
	err = caching.CACHE.Set(fmt.Sprintf("admin_%v", chat.Id), adminJson)
	error_handling.HandleErr(err)
}
