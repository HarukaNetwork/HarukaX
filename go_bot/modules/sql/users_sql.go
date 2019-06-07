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
	PrivChatId int    `sql:",pk"`
	Chat       string `pg:"fk:ChatId"`
	User       int    `ph:"fk:UserId"`
}

func EnsureBotInDb(u *gotgbot.Updater) {
	// usersLock and defer unlock for thread safety
	models := []interface{}{&Users{}, &Chats{}, &ChatMembers{}}
	for _, model := range models {
		_ = SESSION.CreateTable(model, &orm.CreateTableOptions{FKConstraints: true})
	}

	// Insert bot user only if it doesn't exist already
	botUser := &Users{UserId: u.Dispatcher.Bot.Id, UserName: u.Dispatcher.Bot.UserName}
	err := SESSION.Select(botUser)
	if err != nil {
		er := SESSION.Insert(botUser)
		error_handling.HandleErrorGracefully(er)
	}
}

func UpdateUser(userId int, username string, chatId string, chatName string) {
	username = strings.ToLower(username)

	// insert/update user
	user := &Users{UserName: username, UserId: userId}
	err := SESSION.Select(user)
	if err != nil {
		err := SESSION.Insert(user)
		error_handling.HandleErrorGracefully(err)
	} else {
		user.UserName = username
	}

	if chatId == "nil" || chatName == "nil" {
		return
	}

	chat := &Chats{ChatId: string(chatId)}
	err = SESSION.Select(chat)
	if err != nil {
		chat.ChatName = chatName
		err := SESSION.Insert(chat)
		error_handling.HandleErrorGracefully(err)
	} else {
		chat.ChatName = chatName
	}

	member := &ChatMembers{Chat: chat.ChatId, User: user.UserId}
	err = SESSION.Select(member)
	if err != nil {
		if (member.Chat[0] == '-') || (chatName == "nil" && chatId == "nil") {
			err := SESSION.Insert(member)
			error_handling.HandleErrorGracefully(err)
		}
	}
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
	error_handling.HandleErrorGracefully(err)
	return chatMembers
}

func GetAllChats() []Chats {
	var chats []Chats
	err := SESSION.Model(&chats).Select()
	error_handling.HandleErrorGracefully(err)
	return chats
}

func GetUserNumChats(userId int) int {
	count, err := SESSION.Model(new(ChatMembers)).Where("chat_members.user = ?", userId).SelectAndCount()
	error_handling.HandleErrorGracefully(err)
	return count
}

func NumChats() int {
	count, err := SESSION.Model(new(Chats)).SelectAndCount()
	error_handling.HandleErrorGracefully(err)
	return count
}

func NumUsers() int {
	count, err := SESSION.Model(new(Users)).SelectAndCount()
	error_handling.HandleErrorGracefully(err)
	return count
}

func DelUser(userId int) bool {
	user := &Users{UserId: userId}
	err := SESSION.Select(user)
	if err == nil {
		err := SESSION.Delete(user)
		error_handling.HandleErrorGracefully(err)
		return true
	}
	err = SESSION.Delete(&ChatMembers{User: userId})
	return false
}