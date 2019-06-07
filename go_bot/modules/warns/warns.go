package warns

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"github.com/atechnohazard/ginko/go_bot/modules/sql"
	"github.com/atechnohazard/ginko/go_bot/modules/utils/chat_status"
	"github.com/atechnohazard/ginko/go_bot/modules/utils/extraction"
	"github.com/atechnohazard/ginko/go_bot/modules/utils/helpers"
	"html"
	"log"
	"regexp"
	"strconv"
)

func warn(u *ext.User, c *ext.Chat, reason string, m *ext.Message) error {
	var err error
	var reply string

	if chat_status.IsUserAdmin(c, u.Id, nil) {
		_, err = m.ReplyText("Damn admins, can't even be warned!")
		return err
	}

	limit, softWarn := sql.GetWarnSetting(strconv.Itoa(c.Id))
	numWarns, reasons := sql.WarnUser(strconv.Itoa(u.Id), strconv.Itoa(c.Id), reason)

	var keyboard ext.InlineKeyboardMarkup
	if numWarns >= limit {
		sql.ResetWarns(strconv.Itoa(u.Id), strconv.Itoa(c.Id))

		if softWarn {
			_, err = c.UnbanMember(u.Id)
			reply = fmt.Sprintf("%d warnings, %s has been kicked!", limit, helpers.MentionHtml(u.Id, u.FirstName))
			if err != nil {
				return err
			}
		} else {
			_, err = c.KickMember(u.Id)
			reply = fmt.Sprintf("%d warnings, %s has been banned!", limit, helpers.MentionHtml(u.Id, u.FirstName))
			if err != nil {
				return err
			}
		}
		for _, warnReason := range reasons {
			reply += fmt.Sprintf("\n - %v", html.EscapeString(warnReason))
		}
	} else {
		kb := make([][]ext.InlineKeyboardButton, 1)
		kb[0] = make([]ext.InlineKeyboardButton, 1)
		kb[0][0] = ext.InlineKeyboardButton{Text: "Remove warn", CallbackData: fmt.Sprintf("rmWarn(%v)", u.Id)}
		keyboard = ext.InlineKeyboardMarkup{InlineKeyboard: &kb}
		reply = fmt.Sprintf("%v has %v/%v warnings... watch out!", helpers.MentionHtml(u.Id, u.FirstName), numWarns, limit)

		if reason != "" {
			reply += fmt.Sprintf("\nReason for last warn:\n%v", html.EscapeString(reason))
		}
	}

	msg := c.Bot.NewSendableMessage(c.Id, reply)
	msg.ParseMode = parsemode.Html
	msg.ReplyToMessageId = m.MessageId
	msg.ReplyMarkup = &keyboard
	_, err = msg.Send()
	if err != nil {
		msg.ReplyToMessageId = 0
		_, err = msg.Send()
	}
	return err
}

func warnUser(_ ext.Bot, u *gotgbot.Update, args []string) error {
	m := u.EffectiveMessage
	c := u.EffectiveChat

	userId, reason := extraction.ExtractUserAndText(m, args)

	if userId != 0 {
		if m.ReplyToMessage != nil && m.ReplyToMessage.From.Id == userId {
			return warn(m.ReplyToMessage.From, c, reason, m.ReplyToMessage)
		} else {
			chatMember, err := c.GetMember(userId)
			if err != nil {
				return err
			}
			return warn(chatMember.User, c, reason, m)
		}
	} else {
		_, err := m.ReplyText("No user was designated!")
		return err
	}
}

func button(bot ext.Bot, u *gotgbot.Update) error {
	query := u.CallbackQuery
	user := u.EffectiveUser
	pattern, _ := regexp.Compile(`rmWarn\((.+?)\)`)

	if pattern.MatchString(query.Data) {
		userId := pattern.FindStringSubmatch(query.Data)[1]
		chat := u.EffectiveChat
		res := sql.RemoveWarn(userId, strconv.Itoa(chat.Id))
		if res {
			msg := bot.NewSendableEditMessageText(chat.Id, u.EffectiveMessage.MessageId, fmt.Sprintf("Warn removed by %v.", helpers.MentionHtml(user.Id, user.FirstName)))
			msg.ParseMode = parsemode.Html
			_, err := msg.Send()
			return err
		} else {
			_, err := u.EffectiveMessage.EditText("User already has no warns.")
			return err
		}
	}
	return nil
}

func resetWarns(bot ext.Bot, u *gotgbot.Update, args []string) error {
	message := u.EffectiveMessage
	chat := u.EffectiveChat
	userId := extraction.ExtractUser(message, args)

	if userId != 0 {
		sql.ResetWarns(strconv.Itoa(userId), strconv.Itoa(chat.Id))
		_, err := message.ReplyText("Warnings have been reset!")
		return err
	} else {
		_, err := message.ReplyText("No user has been designated!")
		return err
	}
}

func LoadWarns(u *gotgbot.Updater) {
	defer log.Println("Loadin module warns")
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("warn", warnUser))
	u.Dispatcher.AddHandler(handlers.NewCallback("rmWarn", button))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("resetwarns", resetWarns))
}
