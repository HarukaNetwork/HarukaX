package feds

import (
	"fmt"
	"github.com/ATechnoHazard/ginko/go_bot"
	"github.com/ATechnoHazard/ginko/go_bot/modules/sql"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/error_handling"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/extraction"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/helpers"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"log"
	"strconv"
)

func fedBan(bot ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat
	user := u.EffectiveUser
	msg := u.EffectiveMessage
	fedId := sql.GetFedId(strconv.Itoa(chat.Id))

	if fedId == "" {
		_, err := msg.ReplyText("This chat is not part of any federation!")
		return err
	}

	userId, reason := extraction.ExtractUserAndText(msg, args)

	if userId == 0 {
		_, err := msg.ReplyText("Try targeting a user next time bud.")
		return err
	}

	fed := sql.GetFedInfo(fedId)

	if sql.IsUserFedAdmin(fedId, strconv.Itoa(user.Id)) == "" {
		_, err := msg.ReplyHTMLf("You aren't a federation admin for <b>%v</b>", fed.FedName)
		return err
	}

	fbannedUser := sql.GetFbanUser(fedId, strconv.Itoa(userId))

	if strconv.Itoa(userId) == fed.OwnerId {
		_, err := msg.ReplyText("Why are you trying to fban the federation owner?")
		return err
	}

	if sql.IsUserFedAdmin(fedId, strconv.Itoa(userId)) != "" {
		_, err := msg.ReplyText("Why are you trying to fban a federation admin?")
		return err
	}

	if userId == go_bot.BotConfig.OwnerId {
		_, err := msg.ReplyText("I'm not fbanning my owner!")
		return err
	}

	for _, id := range go_bot.BotConfig.SudoUsers {
		sudoId, _ := strconv.Atoi(id)
		if userId == sudoId {
			_, err := msg.ReplyText("I'm not fbanning a sudo user!")
			return err
		}
	}

	if reason == "" {
		reason = "No reason."
	}

	go sql.FbanUser(fedId, strconv.Itoa(userId), reason)
	member, _ := bot.GetChat(userId)

	if fbannedUser == nil {
		_, err := msg.ReplyHTMLf("Beginning federation ban of %v in %v.", helpers.MentionHtml(member.Id, member.FirstName), fed.FedName)
		error_handling.HandleErr(err)
		go func(bot ext.Bot, user *ext.Chat, userId int, federations *sql.Federations, reason string) {
			for _, chat := range sql.AllFedChats(fedId) {
				chatId, err := strconv.Atoi(chat)
				error_handling.HandleErr(err)
				_, err = bot.KickChatMember(chatId, userId)
				error_handling.HandleErr(err)

				_, err = bot.SendMessageHTML(chatId, fmt.Sprintf("User %v is banned in the current federation " +
					"(%v), and so has been removed." +
					"\n<b>Reason</b>: %v", helpers.MentionHtml(member.Id, member.FirstName), fed.FedName, reason))
			}
		}(bot, member, userId, fed, reason)

		_, err = msg.ReplyHTMLf("<b>New FedBan</b>" +
			"\n<b>Fed</b>: %v" +
			"\n<b>FedAdmin</b>: %v" +
			"\n<b>User</b>: %v" +
			"\n<b>User ID</b>: <code>%v</code>" +
			"\n<b>Reason</b>: %v", fed.FedName, helpers.MentionHtml(user.Id, user.FirstName), helpers.MentionHtml(member.Id, member.FirstName),
			member.Id, reason)
		return err
	} else {
		_, err := msg.ReplyHTMLf("<b>FedBan Reason update</b>" +
			"\n<b>Fed</b>: %v" +
			"\n<b>FedAdmin</b>: %v" +
			"\n<b>User</b>: %v" +
			"\n<b>User ID</b>: <code>%v</code>" +
			"\n<b>Previous Reason</b>: %v" +
			"\n<b>New Reason</b>: %v", fed.FedName, helpers.MentionHtml(user.Id, user.FirstName), helpers.MentionHtml(member.Id, member.FirstName),
			member.Id, fbannedUser.Reason, reason)
		return err
	}
}

func LoadFeds(u *gotgbot.Updater) {
	defer log.Println("Loading module feds")
	u.Dispatcher.AddHandler(handlers.NewCommand("newfed", newFed))
	u.Dispatcher.AddHandler(handlers.NewCommand("delfed", delFed))
	u.Dispatcher.AddHandler(handlers.NewCommand("chatfed", chatFed))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("joinfed", joinFed))
	u.Dispatcher.AddHandler(handlers.NewCommand("leavefed", leaveFed))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("fedpromote", fedPromote))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("feddemote", fedDemote))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("fedinfo", fedInfo))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("fedadmins", fedAdmins))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("fban", fedBan))
}
