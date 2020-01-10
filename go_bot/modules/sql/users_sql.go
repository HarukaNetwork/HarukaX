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

package sql

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/wI2L/jettison"

	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/caching"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/error_handling"
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
}

func UpdateUser(userId int, username string, chatId string, chatName string) {
	username = strings.ToLower(username)
	defer func(name string) {
		go cacheUser()
	}(username)
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
}

func GetUserIdByName(username string) *User {
	username = strings.ToLower(username)

	userJson, err := caching.CACHE.Get("users")
	if err != nil {
		go cacheUser()
	}

	log.Println(string(userJson))

	var users []User
	_ = json.Unmarshal(userJson, &users)
	for _, user := range users {
		if user.UserName == username {
			return &user
		}
	}

	return nil
}

func cacheUser() {
	user := &User{}
	var users []User
	SESSION.Model(user).Find(&users)
	userJson, _ := jettison.Marshal(users)
	err := caching.CACHE.Set("users", userJson)
	error_handling.HandleErr(err)
}
