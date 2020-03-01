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

package users

import (
	"log"
	"strconv"

	"github.com/HarukaNetwork/HarukaX/harukax/modules/sql"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/handlers/Filters"
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
