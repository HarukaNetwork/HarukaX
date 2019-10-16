package warns

import (
	"fmt"
	"github.com/ATechnoHazard/ginko/go_bot/modules/sql"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/chat_status"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/extraction"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/helpers"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/handlers/Filters"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
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
		go sql.ResetWarns(strconv.Itoa(u.Id), strconv.Itoa(c.Id))

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
	message := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser
	msg := u.EffectiveMessage

	userId, reason := extraction.ExtractUserAndText(message, args)

	// Check permissions
	if !chat_status.RequireUserAdmin(chat, msg, user.Id, nil) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.RequireBotAdmin(chat, msg) {
		return gotgbot.EndGroups{}
	}

	if userId != 0 {
		if message.ReplyToMessage != nil && message.ReplyToMessage.From.Id == userId {
			return warn(message.ReplyToMessage.From, chat, reason, message.ReplyToMessage)
		} else {
			chatMember, err := chat.GetMember(userId)
			if err != nil {
				return err
			}
			return warn(chatMember.User, chat, reason, message)
		}
	} else {
		_, err := message.ReplyText("No user was designated!")
		return err
	}
}

func button(bot ext.Bot, u *gotgbot.Update) error {
	query := u.CallbackQuery
	user := u.EffectiveUser
	chat := u.EffectiveChat
	pattern, _ := regexp.Compile(`rmWarn\((.+?)\)`)

	// Check permissions
	if !chat_status.IsUserAdmin(chat, user.Id, nil) {
		return gotgbot.EndGroups{}
	}

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

func resetWarns(_ ext.Bot, u *gotgbot.Update, args []string) error {
	message := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser
	msg := u.EffectiveMessage

	userId := extraction.ExtractUser(message, args)

	// Check permissions
	if !chat_status.RequireUserAdmin(chat, msg, user.Id, nil) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.RequireBotAdmin(chat, msg) {
		return gotgbot.EndGroups{}
	}

	if userId != 0 {
		go sql.ResetWarns(strconv.Itoa(userId), strconv.Itoa(chat.Id))
		_, err := message.ReplyText("Warnings have been reset!")
		return err
	} else {
		_, err := message.ReplyText("No user has been designated!")
		return err
	}
}

func warns(_ ext.Bot, u *gotgbot.Update, args []string) error {
	message := u.EffectiveMessage
	chat := u.EffectiveChat
	userId := extraction.ExtractUser(message, args)
	if userId == 0 {
		userId = u.EffectiveUser.Id
	}
	numWarns, reasons := sql.GetWarns(strconv.Itoa(userId), strconv.Itoa(chat.Id))
	text := ""

	if numWarns != 0 {
		limit, _ := sql.GetWarnSetting(strconv.Itoa(chat.Id))
		if len(reasons) > 0 {
			text = fmt.Sprintf("This user has %v/%v warnings, for the following reasons:", numWarns, limit)
			for _, reason := range reasons {
				text += fmt.Sprintf("\n - %v", reason)
			}
			msgs := helpers.SplitMessage(text)
			for _, msg := range msgs {
				_, err := u.EffectiveMessage.ReplyText(msg)
				if err != nil {
					return err
				}
			}
		} else {
			_, err := u.EffectiveMessage.ReplyText(fmt.Sprintf("User has %v/%v warnings, but no reasons for any of them.", numWarns, limit))
			return err
		}
	} else {
		_, err := u.EffectiveMessage.ReplyText("This user hasn't got any warnings!")
		return err
	}
	return nil
}

var TextAndGroupFilter handlers.FilterFunc = func(message *ext.Message) bool {
	return (extraction.ExtractText(message) != "") && (Filters.Group(message))
}

func LoadWarns(u *gotgbot.Updater) {
	defer log.Println("Loading module warns")
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("warn", []rune{'/', '!'}, warnUser))
	u.Dispatcher.AddHandler(handlers.NewCallback("rmWarn", button))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("resetwarns", []rune{'/', '!'}, resetWarns))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("warns", []rune{'/', '!'}, warns))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("addwarn", []rune{'/', '!'}, addWarnFilter))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("nowarn", []rune{'/', '!'}, removeWarnFilter))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("rmwarn", []rune{'/', '!'}, removeWarnFilter)) // Just an alias for nowarn
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("warnlist", []rune{'/', '!'}, listWarnFilters))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("warnlimit", []rune{'/', '!'}, setWarnLimit))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("strongwarn", []rune{'/', '!'}, setWarnStrength))
	u.Dispatcher.AddHandlerToGroup(handlers.NewMessage(TextAndGroupFilter, replyFilter), 9)
}
