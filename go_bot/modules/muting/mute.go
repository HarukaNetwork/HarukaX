package muting

import (
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/chat_status"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/error_handling"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/extraction"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/helpers"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/string_handling"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"log"
	"strings"
)

func mute(bot ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat
	user := u.EffectiveUser
	msg := u.EffectiveMessage

	// Permission checks
	if !chat_status.RequireBotAdmin(chat, msg) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.RequireUserAdmin(chat, msg, user.Id, nil) {
		return gotgbot.EndGroups{}
	}

	userId := extraction.ExtractUser(msg, args)
	if userId == 0 {
		_, err := msg.ReplyText("You'll need to either give me a username to mute, or reply to someone to be muted.")
		return err
	}

	if userId == bot.Id {
		_, err := msg.ReplyText("No u")
		return err
	}

	member, _ := chat.GetMember(userId)

	if member != nil {
		if chat_status.IsUserAdmin(chat, userId, member) {
			_, err := msg.ReplyText("Afraid I can't stop an admin from talking!")
			return err
		} else {
			_, err := bot.RestrictChatMember(chat.Id, userId)
			error_handling.HandleErr(err)
			_, err = msg.ReplyHTMLf("Shush!\nMuted %v.", helpers.MentionHtml(member.User.Id, member.User.FirstName))
			return err
		}
	} else {
		_, err := msg.ReplyText("This user isn't in the chat!")
		return err
	}
}

func unmute(bot ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat
	user := u.EffectiveUser
	msg := u.EffectiveMessage

	// Permission checks
	if !chat_status.RequireBotAdmin(chat, msg) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.RequireUserAdmin(chat, msg, user.Id, nil) {
		return gotgbot.EndGroups{}
	}

	userId := extraction.ExtractUser(msg, args)

	if userId == 0 {
		_, err := msg.ReplyText("You'll need to either give me a username to unmute, or reply to someone to be unmuted.")
		return err
	}

	member, err := chat.GetMember(userId)
	error_handling.HandleErr(err)

	if member != nil {
		if chat_status.IsUserAdmin(chat, userId, member) {
			_, err := msg.ReplyText("This is an admin, what do you expect me to do?")
			return err
		} else {
			_, err := bot.UnRestrictChatMember(chat.Id, userId)
			error_handling.HandleErr(err)
			_, err = msg.ReplyHTMLf("Unmuted %v!", helpers.MentionHtml(userId, member.User.FirstName))
			return err
		}
	} else {
		_, err := msg.ReplyText("This user isn't even in the chat, unmuting them won't make them talk more than they " +
			"already do!")
		return err
	}
}

func tempMute(bot ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat
	user := u.EffectiveUser
	msg := u.EffectiveMessage

	// Permission checks
	if !chat_status.RequireBotAdmin(chat, msg) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.RequireUserAdmin(chat, msg, user.Id, nil) {
		return gotgbot.EndGroups{}
	}

	userId, time := extraction.ExtractUserAndText(msg, args)

	if userId == 0 {
		_, err := msg.ReplyText("You don't seem to be referring to a user.")
		return err
	}

	member, err := chat.GetMember(userId)
	if err != nil {
		if err.Error() == "User not found" {
			_, err := msg.ReplyText("I can't seem to find this user!")
			return err
		}
	}

	if chat_status.IsUserAdmin(chat, userId, member) {
		_, err := msg.ReplyText("I really wish I could mute admins...")
		return err
	}

	if userId == bot.Id {
		_, err := msg.ReplyText("No u")
		return err
	}

	if time == "" {
		_, err := msg.ReplyText("You haven't specified a time to mute this user for!")
		return err
	}

	splitTime := strings.SplitN(time, " ", 2)
	timeVal := strings.ToLower(splitTime[0])

	muteTime := string_handling.ExtractTime(msg, timeVal)

	newMsg := bot.NewSendableRestrictChatMember(chat.Id, userId)
	newMsg.UntilDate = muteTime
	_, err = newMsg.Send()
	if err != nil {
		_, err := msg.ReplyText("Press F, I can't seem to mute this user.")
		error_handling.HandleErr(err)
	}
	_, err = msg.ReplyHTMLf("Muted for <b>%s</b>!", timeVal)
	return err
}

func LoadMuting(u *gotgbot.Updater) {
	defer log.Println("Loading module muting")
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("mute", mute))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("unmute", unmute))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("tmute", tempMute))
}
