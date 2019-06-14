package sql

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/atechnohazard/ginko/go_bot/modules/utils/error_handling"
	"github.com/go-pg/pg/orm"
	"strings"
)

type Users struct {
	UserId   int `sql:",pk"`
	UserName string
}

func (u Users) String() string {
	return fmt.Sprintf("User<%s (%d)>", u.UserName, u.UserId)
}

type Chats struct {
	ChatId   string `sql:",pk"`
	ChatName string
}

func (c Chats) String() string {
	return fmt.Sprintf("<Chat %s (%s)>", c.ChatName, c.ChatId)
}

type ChatMembers struct {
	Chat       string `sql:",pk" pg:"fk:ChatId"`
	User       int    `sql:",pk" pg:"fk:UserId"`
}

func EnsureBotInDb(u *gotgbot.Updater) {
	models := []interface{}{&Users{}, &Chats{}, &ChatMembers{}}
	for _, model := range models {
		_ = SESSION.CreateTable(model, &orm.CreateTableOptions{FKConstraints: true})
	}

	// Insert bot user only if it doesn't exist already
	botUser := &Users{UserId: u.Dispatcher.Bot.Id, UserName: u.Dispatcher.Bot.UserName}
	_, err := SESSION.Model(botUser).OnConflict("(user_id) DO UPDATE").Set("user_name = EXCLUDED.user_name").Insert()
	error_handling.HandleErr(err)
}

func UpdateUser(userId int, username string, chatId string, chatName string) {
	username = strings.ToLower(username)

	// upsert user
	user := &Users{UserName: username, UserId: userId}
	_, err := SESSION.Model(user).OnConflict("(user_id) DO UPDATE").Set("user_name = EXCLUDED.user_name").Insert()
	error_handling.HandleErr(err)

	if chatId == "nil" || chatName == "nil" {
		return
	}

	// upsert chat
	chat := &Chats{ChatId: string(chatId), ChatName: chatName}
	_, err = SESSION.Model(chat).OnConflict("(chat_id) DO UPDATE").Set("chat_name = EXCLUDED.chat_name").Insert()
	error_handling.HandleErr(err)

	// upsert chat_member
	member := &ChatMembers{Chat: chat.ChatId, User: user.UserId}
	_ = SESSION.Insert(member)
}

func GetUserIdByName(username string) *Users {
	username = strings.ToLower(username)
	user := new(Users)
	err := SESSION.Model(user).Where("user_name = ?", username).Select()
	if err != nil {
		return new(Users)
	}
	return user
}

func GetNameByUserId(userId int) *Users {
	user := &Users{UserId: userId}
	err := SESSION.Select(user)
	if err != nil {
		return new(Users)
	}
	return user
}

func GetChatMembers(chatId string) []ChatMembers {
	var chatMembers []ChatMembers
	err := SESSION.Model(&chatMembers).Where("Chat_Members.chat = ?", chatId).Select()
	error_handling.HandleErr(err)
	return chatMembers
}

func GetAllChats() []Chats {
	var chats []Chats
	err := SESSION.Model(&chats).Select()
	error_handling.HandleErr(err)
	return chats
}

func GetUserNumChats(userId int) int {
	count, err := SESSION.Model(new(ChatMembers)).Where("chat_members.user = ?", userId).SelectAndCount()
	error_handling.HandleErr(err)
	return count
}

func NumChats() int {
	count, err := SESSION.Model(new(Chats)).SelectAndCount()
	error_handling.HandleErr(err)
	return count
}

func NumUsers() int {
	count, err := SESSION.Model(new(Users)).SelectAndCount()
	error_handling.HandleErr(err)
	return count
}

func DelUser(userId int) bool {
	user := &Users{UserId: userId}
	err := SESSION.Select(user)
	if err == nil {
		err := SESSION.Delete(user)
		error_handling.HandleErr(err)
		err = SESSION.Delete(&ChatMembers{User: userId})
		error_handling.HandleErr(err)
		return true
	}
	return false
}
