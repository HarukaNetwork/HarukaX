package chat_status

import (
	"github.com/ATechnoHazard/ginko/go_bot"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/error_handling"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"strconv"
)

func CanDelete(chat *ext.Chat, botId int) bool {
	k, err := chat.GetMember(botId)
	error_handling.HandleErr(err)
	return k.CanDeleteMessages
}

func IsUserBanProtected(chat *ext.Chat, userId int, member *ext.ChatMember) bool {
	if chat.Type == "private" || contains(go_bot.BotConfig.SudoUsers, strconv.Itoa(userId)) || chat.AllMembersAdmin {
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

func IsUserAdmin(chat *ext.Chat, userId int, member *ext.ChatMember) bool {
	if chat.Type == "private" || contains(go_bot.BotConfig.SudoUsers, strconv.Itoa(userId)) || chat.AllMembersAdmin {
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

func IsBotAdmin(chat *ext.Chat, member *ext.ChatMember) bool {
	if chat.Type == "private" || chat.AllMembersAdmin {
		return true
	}
	if member == nil {
		mem, err := chat.GetMember(chat.Bot.Id)
		error_handling.HandleErr(err)
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
		_, err := msg.ReplyText("You must be an admin to perform this action.")
		error_handling.HandleErr(err)
		return false
	}
	return true
}

func RequireUserAdmin(chat *ext.Chat, msg *ext.Message, userId int, member *ext.ChatMember) bool {
	if !IsUserAdmin(chat, userId, member) {
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
	if !botChatMember.CanRestrictMembers{
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
