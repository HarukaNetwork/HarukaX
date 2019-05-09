package admin

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"github.com/atechnohazard/ginko/go_bot/modules/utils/chat_status"
	"github.com/atechnohazard/ginko/go_bot/modules/utils/error_handling"
	"github.com/atechnohazard/ginko/go_bot/modules/utils/extraction"
	"github.com/atechnohazard/ginko/go_bot/modules/utils/string_handling"
	"html"
	"log"
	"strconv"
	"strings"
)

func promote(bot ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat
	chatId := chat.Id
	message := u.EffectiveMessage
	user := u.EffectiveUser

	if !chat_status.RequireBotAdmin(chat) {
		return nil
	}
	if !chat_status.RequireUserAdmin(chat, user.Id, nil) {
		return nil
	}
	if !chat_status.CanPromote(bot, chat) {
		return nil
	}

	userId := extraction.ExtractUser(message, args)
	if userId == 0 {
		_, err := message.ReplyText("This user is ded mate.")
		error_handling.HandleErrorGracefully(err)
		return nil
	}

	userMember, err := chat.GetMember(userId)
	error_handling.HandleErrorGracefully(err)

	if userMember.Status == "administrator" || userMember.Status == "creator" {
		_, err := message.ReplyText("Am I supposed to give them a second star or something?")
		error_handling.HandleErrorGracefully(err)
		return nil
	}

	if userId == bot.Id {
		_, err := message.ReplyText("If only I could do this to myself ;_;")
		error_handling.HandleErrorGracefully(err)
		return nil
	}

	botMember, err := chat.GetMember(bot.Id)
	error_handling.HandleErrorGracefully(err)

	sendablePromoteChatMember := bot.NewSendablePromoteChatMember(chatId, userId)
	sendablePromoteChatMember.CanDeleteMessages = botMember.CanDeleteMessages
	sendablePromoteChatMember.CanChangeInfo = botMember.CanDeleteMessages
	sendablePromoteChatMember.CanEditMessages = botMember.CanEditMessages
	sendablePromoteChatMember.CanPostMessages = botMember.CanPostMessages
	sendablePromoteChatMember.CanInviteUsers = botMember.CanInviteUsers
	sendablePromoteChatMember.CanPinMessages = botMember.CanPinMessages
	sendablePromoteChatMember.CanRestrictMembers = botMember.CanRestrictMembers
	sendablePromoteChatMember.CanPromoteMembers = botMember.CanPromoteMembers

	_, err = sendablePromoteChatMember.Send()
	error_handling.HandleErrorGracefully(err)

	_, err = message.ReplyText("Successfully promoted!")
	error_handling.HandleErrorGracefully(err)

	return nil
}

func demote(bot ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat
	chatId := chat.Id
	message := u.EffectiveMessage
	user := u.EffectiveUser

	if !chat_status.RequireBotAdmin(chat) {
		return nil
	}
	if !chat_status.RequireUserAdmin(chat, user.Id, nil) {
		return nil
	}
	if !chat_status.CanPromote(bot, chat) {
		return nil
	}

	userId := extraction.ExtractUser(message, args)
	if userId == 0 {
		_, err := message.ReplyText("This user is ded mate.")
		error_handling.HandleErrorGracefully(err)
		return nil
	}

	userMember, err := chat.GetMember(userId)
	error_handling.HandleErrorGracefully(err)

	if !(userMember.Status == "administrator") {
		_, err := message.ReplyText("Can't demote what wasn't promoted!")
		error_handling.HandleErrorGracefully(err)
		return nil
	}

	if userMember.Status == "creator" {
		_, err := message.ReplyText("This person CREATED the chat, how would I demote them?")
		error_handling.HandleErrorGracefully(err)
		return nil
	}

	if userId == bot.Id {
		_, err := message.ReplyText("Pls no sir ;_;")
		error_handling.HandleErrorGracefully(err)
		return nil
	}

	bb, err := bot.DemoteChatMember(chatId, userId)
	if err != nil || !bb {
		log.Println(err)
		_, err := message.ReplyText("Could not demote. I might not be admin, or the admin status was appointed by another user, so I can't act upon them!")
		error_handling.HandleErrorGracefully(err)
		return nil
	}

	_, err = message.ReplyText("Successfully demoted!")
	error_handling.HandleErrorGracefully(err)


	return nil
}

func pin(bot ext.Bot, u *gotgbot.Update, args []string) error {
	user := u.EffectiveUser
	chat := u.EffectiveChat

	// Check permissions
	if !chat_status.RequireUserAdmin(chat, user.Id, nil) {
		return nil
	}
	if !chat_status.RequireBotAdmin(chat) {
		return nil
	}
	if !chat_status.CanPin(bot, chat) {
		return nil
	}

	isGroup := chat.Type != "private" && chat.Type != "channel"
	prevMessage := u.EffectiveMessage.ReplyToMessage
	isSilent := true

	if len(args) > 0 {
		isSilent = !(strings.ToLower(args[0]) == "notify" || strings.ToLower(args[0]) == "loud" || strings.ToLower(args[0]) == "violent")
	}

	if prevMessage != nil && isGroup {
		sendable := bot.NewSendablePinChatMessage(chat.Id, prevMessage.MessageId)
		sendable.DisableNotification = isSilent
		bb, err := sendable.Send()
		if err != nil || !bb {
			log.Println(err, bb)
			return nil
		}
	}
	return nil
}

func unpin(bot ext.Bot, u *gotgbot.Update) error {
	user := u.EffectiveUser
	chat := u.EffectiveChat

	// Check permissions
	if !chat_status.RequireUserAdmin(chat, user.Id, nil) {
		return nil
	}
	if !chat_status.RequireBotAdmin(chat) {
		return nil
	}
	if !chat_status.CanPin(bot, chat) {
		return nil
	}

	_, err := bot.UnpinChatMessage(chat.Id)
	if err != nil {
		if !(err.Error() == "Bad Request: CHAT_NOT_MODIFIED"){
			return err
		}
	}
	return nil
}

func invitelink(bot ext.Bot, u *gotgbot.Update) error {
	user := u.EffectiveUser
	chat := u.EffectiveChat
	message := u.EffectiveMessage

	// Check permissions
	if !chat_status.RequireUserAdmin(chat, user.Id, nil) {
		return nil
	}
	if !chat_status.RequireBotAdmin(chat) {
		return nil
	}

	if chat.Username != "" {
		_, err := message.ReplyText(chat.Username)
		error_handling.HandleErrorGracefully(err)
		return nil
	} else if chat.Type == "supergroup" || chat.Type == "channel" {
		botMember, err := chat.GetMember(bot.Id)
		error_handling.HandleErrorGracefully(err)
		if botMember.CanInviteUsers {
			inviteLink, err := bot.ExportChatInviteLink(chat.Id)
			error_handling.HandleErrorGracefully(err)
			_, err = message.ReplyText(inviteLink)
			error_handling.HandleErrorGracefully(err)
			return nil
		} else {
			_, err := message.ReplyText("I don't have access to the invite link, try changing my permissions!")
			error_handling.HandleErrorGracefully(err)
			return nil
		}
	} else {
		_, err := message.ReplyText("I can only give you invite links for supergroups and channels, sorry!")
		error_handling.HandleErrorGracefully(err)
		return nil
	}
}

func adminlist(bot ext.Bot, u *gotgbot.Update) error {
	admins, err := u.EffectiveChat.GetAdministrators()
	error_handling.HandleErrorGracefully(err)
	var addendum string
	if u.EffectiveChat.Title != "" {
		addendum = u.EffectiveChat.Title
	} else {
		addendum = "This chat"
	}
	text := fmt.Sprintf("Admins in <b>%s</b>:", addendum)
	for _, admin := range admins {
		user := admin.User
		name := string_handling.FormatText("[{urltext}](tg://user?id={userid})", "{urltext}", user.FirstName + user.LastName, "{userid}", strconv.Itoa(user.Id))
		if user.Username != "" {
			name = html.EscapeString("@" + user.Username)
			text += fmt.Sprintf("\n - %s", name)
		}
	}
	msg := bot.NewSendableMessage(u.EffectiveChat.Id, text)
	msg.ParseMode = parsemode.Html
	msg.ReplyToMessageId = u.EffectiveMessage.MessageId
	_, err = msg.Send()
	error_handling.HandleErrorGracefully(err)
	return nil
}

func LoadAdmin(u *gotgbot.Updater) {
	defer log.Println("Loading module admin")
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("promote", promote))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("demote", demote))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("pin", pin))
	u.Dispatcher.AddHandler(handlers.NewCommand("unpin", unpin))
	u.Dispatcher.AddHandler(handlers.NewCommand("invitelink", invitelink))
	u.Dispatcher.AddHandler(handlers.NewCommand("adminlist", adminlist))
}
