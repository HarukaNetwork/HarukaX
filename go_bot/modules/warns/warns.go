package warns

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/handlers/Filters"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"github.com/atechnohazard/ginko/go_bot/modules/sql"
	"github.com/atechnohazard/ginko/go_bot/modules/utils/chat_status"
	"github.com/atechnohazard/ginko/go_bot/modules/utils/error_handling"
	"github.com/atechnohazard/ginko/go_bot/modules/utils/extraction"
	"github.com/atechnohazard/ginko/go_bot/modules/utils/helpers"
	"html"
	"log"
	"regexp"
	"strconv"
	"strings"
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
	message := u.EffectiveMessage
	chat := u.EffectiveChat
	user := u.EffectiveUser

	userId, reason := extraction.ExtractUserAndText(message, args)

	// Check permissions
	if !chat_status.RequireUserAdmin(chat, user.Id, nil) {
		return nil
	}
	if !chat_status.RequireBotAdmin(chat) {
		return nil
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
		return nil
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
	userId := extraction.ExtractUser(message, args)

	// Check permissions
	if !chat_status.RequireUserAdmin(chat, user.Id, nil) {
		return nil
	}
	if !chat_status.RequireBotAdmin(chat) {
		return nil
	}

	if userId != 0 {
		sql.ResetWarns(strconv.Itoa(userId), strconv.Itoa(chat.Id))
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
	numWarns, reasons := sql.GetWarns(strconv.Itoa(userId), strconv.Itoa(chat.Id))
	text := ""

	if numWarns != 0 {
		limit, _ := sql.GetWarnSetting(strconv.Itoa(chat.Id))
		if reasons != nil {
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

func addWarnFilter(_ ext.Bot, u *gotgbot.Update) error {
	chat := u.EffectiveChat
	msg := u.EffectiveMessage
	user := u.EffectiveUser
	var keyword string
	var content string

	// Check permissions
	if !chat_status.RequireUserAdmin(chat, user.Id, nil) {
		return nil
	}
	if !chat_status.RequireBotAdmin(chat) {
		return nil
	}

	args := strings.SplitN(msg.Text, " ", 2)

	if len(args) < 2 {
		return nil
	}

	extracted := helpers.SplitQuotes(args[1])

	if len(extracted) >= 2 {
		keyword = strings.ToLower(extracted[0])
		content = extracted[1]
	} else {
		return nil
	}

	sql.AddWarnFilter(strconv.Itoa(chat.Id), keyword, content)
	_, err := msg.ReplyText(fmt.Sprintf("Warn handler added for '%v'!", keyword))
	error_handling.HandleErr(err)
	return gotgbot.EndGroups{}
}

func removeWarnFilter(_ ext.Bot, u *gotgbot.Update) error {
	chat := u.EffectiveChat
	user := u.EffectiveUser
	msg := u.EffectiveMessage

	// Check permissions
	if !chat_status.RequireUserAdmin(chat, user.Id, nil) {
		return nil
	}
	if !chat_status.RequireBotAdmin(chat) {
		return nil
	}

	args := strings.SplitN(msg.Text, " ", 2)

	if len(args) < 2 {
		return nil
	}

	extracted := helpers.SplitQuotes(args[1])

	if len(extracted) < 1 {
		return nil
	}

	toRemove := extracted[0]

	chatFilters := sql.GetChatWarnTriggers(strconv.Itoa(chat.Id))

	if chatFilters == nil {
		_, err := msg.ReplyText("No warning filters are active here!")
		return err
	}

	for _, filt := range chatFilters {
		if filt.Keyword == toRemove {
			sql.RemoveWarnFilter(strconv.Itoa(chat.Id), toRemove)
			_, err := msg.ReplyText("Yep, I'll stop warning people for that.")
			error_handling.HandleErr(err)
			return gotgbot.EndGroups{}
		}
	}
	_, err := msg.ReplyText("That's not a current warning filter - run /warnlist for all active warning filters.")
	return err
}

func listWarnFilters(_ ext.Bot, u *gotgbot.Update) error {
	chat := u.EffectiveChat
	allHandlers := sql.GetChatWarnTriggers(strconv.Itoa(chat.Id))

	if allHandlers == nil {
		_, err := u.EffectiveMessage.ReplyText("No warning filters are active here!")
		return err
	}

	filterList := "<b>Current warning filters in this chat:</b>\n"
	for _, handler := range allHandlers {
		entry := fmt.Sprintf(" - %v\n", html.EscapeString(handler.Keyword))
		if len(entry) + len(filterList) > 4096 {
			_, err := u.EffectiveMessage.ReplyHTML(filterList)
			error_handling.HandleErr(err)
			filterList = entry
		} else {
			filterList += entry
		}
	}
	if filterList != "<b>Current warning filters in this chat:</b>\n" {
		_, err := u.EffectiveMessage.ReplyHTML(filterList)
		return err
	}
	return nil
}

func replyFilter(_ ext.Bot, u *gotgbot.Update) error {
	chat := u.EffectiveChat
	message := u.EffectiveMessage

	if !chat_status.RequireBotAdmin(chat) {
		return nil
	}

	chatWarnFilters := sql.GetChatWarnTriggers(strconv.Itoa(chat.Id))
	toMatch := extraction.ExtractText(message)
	if toMatch == "" {
		return nil
	}

	for _, handler := range chatWarnFilters {
		pattern := `( |^|[^\w])` + regexp.QuoteMeta(handler.Keyword) + `( |$|[^\w])`
		re, err := regexp.Compile("(?i)" + pattern)
		error_handling.HandleErr(err)

		if re.MatchString(toMatch) {
			user := u.EffectiveUser
			warnFilter := sql.GetWarnFilter(strconv.Itoa(chat.Id), handler.Keyword)
			return warn(user, chat, warnFilter.Reply, message)
		}
	}
	return gotgbot.ContinueGroups{}
}

var TextAndGroupFilter handlers.FilterFunc = func (message *ext.Message) bool {
	return (extraction.ExtractText(message) != "") && (Filters.Group(message))
}

func LoadWarns(u *gotgbot.Updater) {
	defer log.Println("Loading module warns")
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("warn", warnUser))
	u.Dispatcher.AddHandler(handlers.NewCallback("rmWarn", button))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("resetwarns", resetWarns))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("warns", warns))
	u.Dispatcher.AddHandler(handlers.NewCommand("addwarn", addWarnFilter))
	u.Dispatcher.AddHandler(handlers.NewCommand("nowarn", removeWarnFilter))
	u.Dispatcher.AddHandler(handlers.NewCommand("warnlist", listWarnFilters))
	u.Dispatcher.AddHandlerToGroup(handlers.NewMessage(TextAndGroupFilter, replyFilter), 9)
}
