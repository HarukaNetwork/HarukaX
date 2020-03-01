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

package sql

import (
	"encoding/json"
	"strings"

	"github.com/wI2L/jettison"

	"github.com/HarukaNetwork/HarukaX/harukax/modules/utils/caching"
	"github.com/HarukaNetwork/HarukaX/harukax/modules/utils/error_handling"
	"github.com/PaulSonOfLars/gotgbot"
)

type User struct {
	UserId   int    `gorm:"primary_key" json:"user_id"`
	UserName string `json:"user_name"`
}

type Chat struct {
	ChatId   string `gorm:"primary_key" json:"chat_id"`
	ChatName string `json:"chat_name"`
}

func EnsureBotInDb(u *gotgbot.Updater) {
	// Insert bot user only if it doesn't exist already
	botUser := &User{UserId: u.Dispatcher.Bot.Id, UserName: u.Dispatcher.Bot.UserName}
	SESSION.Save(botUser)
	cacheUser()
}

func UpdateUser(userId int, username string, chatId string, chatName string) {
	username = strings.ToLower(username)
	tx := SESSION.Begin()

	// upsert user
	user := &User{UserId: userId, UserName: username}
	tx.Save(user)

	if chatId == "nil" || chatName == "nil" {
		tx.Commit()
		return
	}

	// upsert chat
	chat := &Chat{ChatId: chatId, ChatName: chatName}
	tx.Save(chat)
	tx.Commit()
	cacheUser()
}

func GetUserIdByName(username string) *User {
	username = strings.ToLower(username)

	userJson, err := caching.CACHE.Get("users")
	var users []User
	if err != nil {
		users = cacheUser()
	}

	_ = json.Unmarshal(userJson, &users)

	for _, user := range users {
		if user.UserName == username {
			return &user
		}
	}

	return nil
}

func cacheUser() []User {
	var users []User
	SESSION.Model(&User{}).Find(&users)
	userJson, _ := jettison.Marshal(users)
	err := caching.CACHE.Set("users", userJson)
	error_handling.HandleErr(err)
	return users
}
