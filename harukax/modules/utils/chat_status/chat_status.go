/*
 *    Copyright © 2020 Haruka Network Development
 *    This file is part of Haruka X.
 *
 *    Haruka X is free software: you can redistribute it and/or modify
 *    it under the terms of the Raphielscape Public License as published by
 *    the Devscapes Open Source Holding GmbH., version 1.d
 *
 *    Haruka X is distributed in the hope that it will be useful,
 *    but WITHOUT ANY WARRANTY; without even the implied warranty of
 *    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *    Devscapes Raphielscape Public License for more details.
 *
 *    You should have received a copy of the Devscapes Raphielscape Public License
 */

package chat_status

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/wI2L/jettison"

	"github.com/HarukaNetwork/HarukaX/harukax"
	"github.com/HarukaNetwork/HarukaX/harukax/modules/utils/caching"
	"github.com/HarukaNetwork/HarukaX/harukax/modules/utils/error_handling"
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
	if chat.Type == "private" || containsID(harukax.BotConfig.SudoUsers, userId) {
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
	if chat.Type == "private" || containsID(harukax.BotConfig.SudoUsers, userId) {
		return true
	}

	var adminList []ext.ChatMember

	// try to find admins in cache
	admins, err := caching.CACHE.Get(fmt.Sprintf("admin_%v", chat.Id))
	if err != nil { // cache miss
		adminList, err = chat.GetAdministrators()
		if err != nil { // not found
			return false
		}
		// found using API, save to cache
		cacheAdmins(adminList, chat.Id)
		return contains(adminList, userId)
	}

	_ = json.Unmarshal(admins, &adminList)
	return contains(adminList, userId)
}

func containsID(sudos []string, userID int) bool {
	for _, a := range sudos {
		if strconv.Itoa(userID) == a {
			return true
		}
	}
	return false
}

func cacheAdmins(adminList []ext.ChatMember, chatID int) {
	adminJson, _ := jettison.Marshal(adminList)
	_ = caching.CACHE.Set(fmt.Sprintf("admin_%v", chatID), adminJson)
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

func RequireUserAdmin(chat *ext.Chat, msg *ext.Message, userId int) bool {
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

func contains(s []ext.ChatMember, e int) bool {
	for _, a := range s {
		if a.User.Id == e {
			return true
		}
	}
	return false
}
