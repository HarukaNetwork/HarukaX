package chat_status

import (
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/atechnohazard/ginko/go_bot"
	"github.com/atechnohazard/ginko/go_bot/modules/utils/error_handling"
	"strconv"
)

func CanDelete(chat *ext.Chat, botId int) bool {
	k, err := chat.GetMember(botId)
	error_handling.HandleErrorGracefully(err)
	return k.CanDeleteMessages
}

func IsUserBanProtected(chat *ext.Chat, userId int, member *ext.ChatMember) bool {
	if chat.Type == "private" || contains(go_bot.BotConfig.SudoUsers, strconv.Itoa(userId)) || chat.AllMembersAdmin {
		return true
	}
	if member == nil {
		mem, err := chat.GetMember(userId)
		error_handling.HandleErrorGracefully(err)
		member = mem
	}
	if member.Status == "administrator" || member.Status == "creator" {
		return true
	} else {
		return false
	}
}

func IsUserAdmin(chat *ext.Chat, userId int, member *ext.ChatMember) bool {
	if chat.Type == "private" || contains(go_bot.BotConfig.SudoUsers, strconv.Itoa(userId)) || chat.AllMembersAdmin {
		return true
	}
	if member == nil {
		mem, err := chat.GetMember(userId)
		error_handling.HandleErrorGracefully(err)
		member = mem
	}
	if member.Status == "administrator" || member.Status == "creator" {
		return false
	} else {
		return true
	}
}

func IsBotAdmin(chat *ext.Chat, member *ext.ChatMember) bool {
	if chat.Type == "private" || chat.AllMembersAdmin {
		return true
	}
	if member == nil {
		mem, err := chat.GetMember(chat.Bot.Id)
		error_handling.HandleErrorGracefully(err)
		member = mem
	}
	if member.Status == "administrator" || member.Status == "creator" {
		return true
	} else {
		return false
	}
}

func RequireBotAdmin(chat *ext.Chat) bool {
	if !IsBotAdmin(chat, nil) {
		_, err := chat.Bot.SendMessage(chat.Id, "I'm not admin!")
		error_handling.HandleErrorGracefully(err)
		return false
	}
	return true
}

func RequireUserAdmin(chat *ext.Chat, userId int, member *ext.ChatMember) bool {
	if !IsUserAdmin(chat, userId, member) {
		_, err := chat.Bot.SendMessage(chat.Id, "You must be an admin to perform this action.")
		error_handling.HandleErrorGracefully(err)
		return false
	}
	return true
}

func IsUserInChat(chat *ext.Chat, userId int) bool {
	member, err := chat.GetMember(userId)
	error_handling.HandleErrorGracefully(err)
	if member.Status == "left" || member.Status == "kicked" {
		return false
	} else {
		return true
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
