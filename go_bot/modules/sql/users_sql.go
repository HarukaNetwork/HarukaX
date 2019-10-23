package sql

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"strings"
)

type User struct {
	UserId   int `gorm:"primary_key"`
	UserName string
}

func (u User) String() string {
	return fmt.Sprintf("User<%s (%d)>", u.UserName, u.UserId)
}

type Chat struct {
	ChatId   string `gorm:"primary_key"`
	ChatName string
}

func (c Chat) String() string {
	return fmt.Sprintf("<Chat %s (%s)>", c.ChatName, c.ChatId)
}

func EnsureBotInDb(u *gotgbot.Updater) {
	// Insert bot user only if it doesn't exist already
	botUser := &User{UserId: u.Dispatcher.Bot.Id, UserName: u.Dispatcher.Bot.UserName}
	SESSION.Save(botUser)
}

func UpdateUser(userId int, username string, chatId string, chatName string) {
	username = strings.ToLower(username)
	tx := SESSION.Begin()

	// upsert user
	user := &User{UserId: userId, UserName: username}
	tx.Save(user)

	if chatId == "nil" || chatName == "nil" {
		return
	}

	// upsert chat
	chat := &Chat{ChatId: chatId, ChatName: chatName}
	tx.Save(chat)
	tx.Commit()
}

func GetUserIdByName(username string) *User {
	username = strings.ToLower(username)
	user := new(User)
	SESSION.Where("user_name = ?", username).First(user)
	return user
}
