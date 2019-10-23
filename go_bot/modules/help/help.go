package help

import (
	"fmt"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/error_handling"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"log"
	"regexp"
)

var markup ext.InlineKeyboardMarkup

func initHelpButtons() {
	helpButtons := [][]ext.InlineKeyboardButton{make([]ext.InlineKeyboardButton, 2), make([]ext.InlineKeyboardButton, 2),
		make([]ext.InlineKeyboardButton, 2), make([]ext.InlineKeyboardButton, 2), make([]ext.InlineKeyboardButton, 2)}

	// First column
	helpButtons[0][0] = ext.InlineKeyboardButton{
		Text:         "admin",
		CallbackData: fmt.Sprintf("help(%v)", "admin"),
	}
	helpButtons[1][0] = ext.InlineKeyboardButton{
		Text:         "bans",
		CallbackData: fmt.Sprintf("help(%v)", "bans"),
	}
	helpButtons[2][0] = ext.InlineKeyboardButton{
		Text:         "blacklist",
		CallbackData: fmt.Sprintf("help(%v)", "blacklist"),
	}
	helpButtons[3][0] = ext.InlineKeyboardButton{
		Text:         "deleting",
		CallbackData: fmt.Sprintf("help(%v)", "deleting"),
	}
	helpButtons[4][0] = ext.InlineKeyboardButton{
		Text:         "federations",
		CallbackData: fmt.Sprintf("help(%v)", "feds"),
	}

	// Second column
	helpButtons[0][1] = ext.InlineKeyboardButton{
		Text:         "misc",
		CallbackData: fmt.Sprintf("help(%v)", "misc"),
	}
	helpButtons[1][1] = ext.InlineKeyboardButton{
		Text:         "muting",
		CallbackData: fmt.Sprintf("help(%v)", "muting"),
	}
	helpButtons[2][1] = ext.InlineKeyboardButton{
		Text:         "notes",
		CallbackData: fmt.Sprintf("help(%v)", "notes"),
	}
	helpButtons[3][1] = ext.InlineKeyboardButton{
		Text:         "users",
		CallbackData: fmt.Sprintf("help(%v)", "users"),
	}
	helpButtons[4][1] = ext.InlineKeyboardButton{
		Text:         "warns",
		CallbackData: fmt.Sprintf("help(%v)", "warns"),
	}

	markup = ext.InlineKeyboardMarkup{InlineKeyboard: &helpButtons}
}

func help(b ext.Bot, u *gotgbot.Update) error {
	msg := b.NewSendableMessage(u.EffectiveChat.Id, "Hey there! I'm Ginko, a group management bot written in Go."+
		"I have a ton of useful features like notes, filters and even a warn system.\n\n"+
		"Commands are preceded with a slash  (/) or an exclamation mark (!)\n\n"+
		"Some basic commands:\n\n"+
		"- /start: duh, you already know what this does\n\n"+
		"- /help: for info on how to use me\n\n"+
		"- /donate: info on who made me and how you can support them\n\n\n"+
		"If you have any bugs reports, questions or suggestions you can head over to @gobotsupport.\n\n"+
		"Have fun using me!")
	msg.ParseMode = parsemode.Html
	msg.ReplyToMessageId = u.EffectiveMessage.MessageId
	msg.ReplyMarkup = &markup
	_, err := msg.Send()
	if err != nil {
		msg.ReplyToMessageId = 0
		_, err = msg.Send()
	}
	return err
}

func buttonHandler(b ext.Bot, u *gotgbot.Update) error {
	query := u.CallbackQuery
	pattern, _ := regexp.Compile(`help\((.+?)\)`)

	if pattern.MatchString(query.Data) {
		module := pattern.FindStringSubmatch(query.Data)[1]
		chat := u.EffectiveChat
		msg := b.NewSendableEditMessageText(chat.Id, u.EffectiveMessage.MessageId, "placeholder")
		msg.ParseMode = parsemode.Html
		backButton := [][]ext.InlineKeyboardButton{{ext.InlineKeyboardButton{
			Text:         "back",
			CallbackData: "help(back)",
		}}}
		backKeyboard := ext.InlineKeyboardMarkup{InlineKeyboard: &backButton}
		msg.ReplyMarkup = &backKeyboard

		switch module {
		case "admin":
			msg.Text = "Here is the help for the <b>Admin</b> module:\n\n" +
				"- /adminlist: list of admins in the chat\n\n" +
				"<b>Admin only:</b>" +
				"- /pin: silently pins the message replied to - add 'loud' or 'notify' to give notifs to users.\n" +
				"- /unpin: unpins the currently pinned message\n" +
				"- /invitelink: gets invitelink\n" +
				"- /promote: promotes the user replied to\n" +
				"- /demote: demotes the user replied to\n"
			break
		case "bans":
			msg.Text = "Here is the help for the <b>Bans</b> module:\n\n" +
				" - /kickme: kicks the user who issued the command\n\n" +
				"<b>Admin only</b>:\n" +
				" - /ban <userhandle>: bans a user. (via handle, or reply)\n" +
				" - /tban <userhandle> x(m/h/d): bans a user for x time. (via handle, or reply). m = minutes, h = hours," +
				" d = days.\n" +
				"- /unban <userhandle>: unbans a user. (via handle, or reply)" +
				" - /kick <userhandle>: kicks a user, (via handle, or reply)"
			break
		case "blacklist":
			break
		case "deleting":
			break
		case "feds":
			break
		case "misc":
			break
		case "muting":
			break
		case "notes":
			break
		case "users":
			break
		case "warns":
			break
		case "back":
			msg.Text = "Hey there! I'm Ginko, a group management bot written in Go." +
				"I have a ton of useful features like notes, filters and even a warn system.\n\n" +
				"Commands are preceded with a slash (/) or an exclamation mark (!)\n\n" +
				"Some basic commands:\n\n" +
				"- /start: duh, you already know what this does\n\n" +
				"- /help: for info on how to use me\n\n" +
				"- /donate: info on who made me and how you can support them\n\n\n" +
				"If you have any bugs reports, questions or suggestions you can head over to @gobotsupport.\n\n" +
				"Have fun using me!"
			msg.ReplyMarkup = &markup
			break
		}

		_, err := msg.Send()
		error_handling.HandleErr(err)
		_, err = b.AnswerCallbackQuery(query.Id)
		return err
	}
	return nil
}

func LoadHelp(u *gotgbot.Updater) {
	defer log.Println("Loading module help")
	initHelpButtons()
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("help", []rune{'/', '!'}, help))
	u.Dispatcher.AddHandler(handlers.NewCallback("help", buttonHandler))
}
