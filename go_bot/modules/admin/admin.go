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

package admin

import (
	"fmt"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/chat_status"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/error_handling"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/extraction"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/helpers"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/string_handling"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
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

	// permission checks
	if chat.Type == "private" {
		_, err := message.ReplyText("This command is meant to be used in a group!")
		return err
	}
	if !chat_status.RequireBotAdmin(chat, message) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.RequireUserAdmin(chat, message, user.Id, nil) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.CanPromote(bot, chat) {
		return gotgbot.EndGroups{}
	}

	userId := extraction.ExtractUser(message, args)
	if userId == 0 {
		_, err := message.ReplyText("This user is ded mate.")
		error_handling.HandleErr(err)
		return nil
	}

	userMember, err := chat.GetMember(userId)
	error_handling.HandleErr(err)

	if userId == bot.Id {
		_, err := message.ReplyText("If only I could do this to myself ;_;")
		error_handling.HandleErr(err)
		return nil
	}

	botMember, err := chat.GetMember(bot.Id)
	error_handling.HandleErr(err)

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
	error_handling.HandleErr(err)

	_, err = message.ReplyHTMLf("Successfully promoted %v!", helpers.MentionHtml(userId, userMember.User.FirstName))

	return err
}

func demote(bot ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat
	chatId := chat.Id
	message := u.EffectiveMessage
	user := u.EffectiveUser

	// permission checks
	if chat.Type == "private" {
		_, err := message.ReplyText("This command is meant to be used in a group!")
		return err
	}
	if !chat_status.RequireBotAdmin(chat, message) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.RequireUserAdmin(chat, message, user.Id, nil) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.CanPromote(bot, chat) {
		return gotgbot.EndGroups{}
	}

	userId := extraction.ExtractUser(message, args)
	if userId == 0 {
		_, err := message.ReplyText("This user is ded mate.")
		return err
	}

	userMember, err := chat.GetMember(userId)
	error_handling.HandleErr(err)

	if userMember.Status == "creator" {
		_, err := message.ReplyText("This person CREATED the chat, how would I demote them?")
		return err
	}

	if userId == bot.Id {
		_, err := message.ReplyText("Pls no sir ;_;")
		return err
	}

	bb, err := bot.DemoteChatMember(chatId, userId)
	if err != nil || !bb {
		log.Println(err)
		_, err := message.ReplyText("Could not demote. I might not be admin, or the admin status was appointed by another user, so I can't act upon them!")
		return err
	}

	_, err = message.ReplyText("Successfully demoted!")
	return err
}

func pin(bot ext.Bot, u *gotgbot.Update, args []string) error {
	user := u.EffectiveUser
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	// Check permissions
	if chat.Type == "private" {
		_, err := msg.ReplyText("This command is meant to be used in a group!")
		return err
	}
	if !chat_status.RequireUserAdmin(chat, msg, user.Id, nil) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.RequireBotAdmin(chat, msg) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.CanPin(bot, chat) {
		return gotgbot.EndGroups{}
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
		_, err := sendable.Send()
		return err
	}
	return nil
}

func unpin(bot ext.Bot, u *gotgbot.Update) error {
	user := u.EffectiveUser
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	// Check permissions
	if chat.Type == "private" {
		_, err := msg.ReplyText("This command is meant to be used in a group!")
		return err
	}
	if !chat_status.RequireUserAdmin(chat, msg, user.Id, nil) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.RequireBotAdmin(chat, msg) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.CanPin(bot, chat) {
		return gotgbot.EndGroups{}
	}

	_, err := bot.UnpinChatMessage(chat.Id)
	return err
}

func invitelink(bot ext.Bot, u *gotgbot.Update) error {
	user := u.EffectiveUser
	chat := u.EffectiveChat
	message := u.EffectiveMessage

	// Check permissions
	if chat.Type == "private" {
		_, err := message.ReplyText("This command is meant to be used in a group!")
		return err
	}
	if !chat_status.RequireUserAdmin(chat, message, user.Id, nil) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.RequireBotAdmin(chat, message) {
		return gotgbot.EndGroups{}
	}

	if chat.Username != "" {
		_, err := message.ReplyText(chat.Username)
		return err
	} else if chat.Type == "supergroup" || chat.Type == "channel" {
		botMember, err := chat.GetMember(bot.Id)
		error_handling.HandleErr(err)
		if botMember.CanInviteUsers {
			inviteLink, err := bot.ExportChatInviteLink(chat.Id)
			error_handling.HandleErr(err)
			_, err = message.ReplyText(inviteLink)
			return err
		} else {
			_, err := message.ReplyText("I don't have access to the invite link, try changing my permissions!")
			return err
		}
	} else {
		_, err := message.ReplyText("I can only give you invite links for supergroups and channels, sorry!")
		return err
	}
}

func adminlist(_ ext.Bot, u *gotgbot.Update) error {
	if u.EffectiveChat.Type == "private" {
		_, err := u.EffectiveMessage.ReplyText("This command is meant to be used in a group!")
		return err
	}

	admins, err := u.EffectiveChat.GetAdministrators()
	error_handling.HandleErr(err)
	var addendum string
	if u.EffectiveChat.Title != "" {
		addendum = u.EffectiveChat.Title
	} else {
		addendum = "This chat"
	}
	text := fmt.Sprintf("Admins in <b>%s</b>:", addendum)
	for _, admin := range admins {
		user := admin.User
		name := string_handling.FormatText("[{urltext}](tg://user?id={userid})", "{urltext}",
			user.FirstName+user.LastName, "{userid}", strconv.Itoa(user.Id))

		if user.Username != "" {
			name = html.EscapeString("@" + user.Username)
			text += fmt.Sprintf("\n - %s", name)
		}
	}
	_, err = u.EffectiveMessage.ReplyHTML(text)
	return err
}

func LoadAdmin(u *gotgbot.Updater) {
	defer log.Println("Loading module admin")
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("promote", []rune{'/', '!'}, promote))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("demote", []rune{'/', '!'}, demote))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("pin", []rune{'/', '!'}, pin))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("unpin", []rune{'/', '!'}, unpin))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("invitelink", []rune{'/', '!'}, invitelink))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("adminlist", []rune{'/', '!'}, adminlist))
}
