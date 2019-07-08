package notes

import (
	"fmt"
	"github.com/ATechnoHazard/ginko/go_bot/modules/sql"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/chat_status"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/error_handling"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/helpers"
	tgmd2html "github.com/PaulSonOfLars/gotg_md2html"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"html"
	"log"
	"strconv"
	"strings"
)

func get(bot ext.Bot, u *gotgbot.Update, noteName string, showNone bool, noFormat bool) error {
	chatId := u.EffectiveChat.Id
	note := sql.GetNote(strconv.Itoa(chatId), noteName)
	msg := u.EffectiveMessage

	replyId := msg.MessageId

	if note != nil {
		if msg.ReplyToMessage != nil {
			replyId = msg.ReplyToMessage.MessageId
		}

		if note.IsReply {
			msgId, _ := strconv.Atoi(note.Value)
			_, err := bot.ForwardMessage(chatId, chatId, msgId)
			if err != nil {
				_, err := msg.ReplyText("Looks like the original sender of this note has deleted " +
					"their message - sorry! I'll remove this note from " +
					"your saved notes.")
				sql.RmNote(strconv.Itoa(chatId), noteName)
				return err
			}
		} else {
			text := note.Value
			keyb := make([][]ext.InlineKeyboardButton, 0)
			buttons := sql.GetButtons(strconv.Itoa(chatId), noteName)
			parseMode := parsemode.Markdown
			btns := make([]tgmd2html.Button, len(buttons))
			for i, btn := range buttons {
				btns[i] = tgmd2html.Button{Name: btn.Name, Content: btn.Url, SameLine: btn.SameLine}
			}

			if noFormat {
				text = tgmd2html.Reverse(note.Value, btns)
				parseMode = ""
			} else {
				keyb = helpers.BuildKeyboard(buttons)
			}

			keyboard := &ext.InlineKeyboardMarkup{InlineKeyboard: &keyb}

			if note.Msgtype == sql.BUTTON_TEXT || note.Msgtype == sql.TEXT {
				msg := bot.NewSendableMessage(chatId, text)
				msg.ParseMode = parseMode
				msg.ReplyMarkup = keyboard
				msg.DisableWebPreview = true
				msg.ReplyToMessageId = replyId
				_, err := msg.Send()
				return err
			} else {

				var err error
				switch note.Msgtype {
				case sql.STICKER:
					msg := bot.NewSendableSticker(chatId)
					msg.ReplyToMessageId = replyId
					msg.FileId = note.File
					msg.ReplyMarkup = keyboard
					_, err = msg.Send()
					break
				case sql.DOCUMENT:
					msg := bot.NewSendableDocument(chatId, text)
					msg.ReplyToMessageId = replyId
					msg.FileId = note.File
					msg.ReplyMarkup = keyboard
					_, err = msg.Send()
					break
				case sql.PHOTO:
					msg := bot.NewSendablePhoto(chatId, text)
					msg.ParseMode = parseMode
					msg.ReplyToMessageId = replyId
					msg.FileId = note.File
					msg.ReplyMarkup = keyboard
					_, err = msg.Send()
					break
				case sql.AUDIO:
					msg := bot.NewSendableAudio(chatId, text)
					msg.ParseMode = parseMode
					msg.ReplyToMessageId = replyId
					msg.FileId = note.File
					msg.ReplyMarkup = keyboard
					_, err = msg.Send()
					break
				case sql.VOICE:
					msg := bot.NewSendableVoice(chatId, text)
					msg.ParseMode = parseMode
					msg.ReplyToMessageId = replyId
					msg.FileId = note.File
					msg.ReplyMarkup = keyboard
					_, err = msg.Send()
					break
				case sql.VIDEO:
					msg := bot.NewSendableVideo(chatId, text)
					msg.ParseMode = parseMode
					msg.ReplyToMessageId = replyId
					msg.FileId = note.File
					msg.ReplyMarkup = keyboard
					_, err = msg.Send()
					break
				}

				if err != nil {
					if err.Error() == "Bad Request: Entity_mention_user_invalid" {
						_, _ = msg.ReplyText("Looks like you tried to mention someone I've never seen before. If you really " +
							"want to mention them, forward one of their messages to me, and I'll be able " +
							"to tag them!")
						return nil
					} else {
						_, _ = msg.ReplyText("This note could not be sent, as it is incorrectly formatted. Ask in " +
							"@GoBotSupport if you can't figure out why!")
						return nil
					}
				}

			}
		}
	} else if showNone {
		_, err := msg.ReplyText("This note doesn't exist!")
		return err
	}
	return nil
}

func cmdGet(bot ext.Bot, u *gotgbot.Update, args []string) error {
	if len(args) >= 2 && strings.ToLower(args[1]) == "noformat" {
		return get(bot, u, args[0], true, true)
	} else if len(args) >= 1 {
		return get(bot, u, args[0], true, false)
	} else {
		_, err := u.EditedMessage.ReplyText("Get rekt")
		return err
	}
}

func hashGet(bot ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage
	fstWord := strings.Split(msg.Text, " ")[0]
	noHash := fstWord[1:]
	return get(bot, u, noHash, false, false)
}

func save(_ ext.Bot, u *gotgbot.Update) error {
	chatId := u.EffectiveChat.Id
	msg := u.EffectiveMessage
	noteName, text, dataType, content, buttons := helpers.GetNoteType(msg)

	if !chat_status.RequireUserAdmin(u.EffectiveChat, msg, u.EffectiveUser.Id, nil) {
		return gotgbot.ContinueGroups{}
	}

	if dataType == -1 {
		_, err := msg.ReplyText("Dude, there's no note!")
		return err
	}

	if len(strings.TrimSpace(text)) == 0 {
		text = noteName
	}

	btns := make([]sql.Button, len(buttons))

	for i, btn := range buttons {
		btns[i] = sql.Button{ChatId: strconv.Itoa(chatId), NoteName: noteName, Name: btn.Name, Url: btn.Content, SameLine: btn.SameLine}
	}

	go sql.AddNoteToDb(strconv.Itoa(chatId), noteName, text, dataType, btns, content)
	_, err := msg.ReplyHTMLf("Added %v!\nGet it with /get %v, or #%v", noteName, noteName, noteName)
	return err
}

func clear(_ ext.Bot, u *gotgbot.Update, args []string) error {
	chatId := u.EffectiveChat.Id

	if !chat_status.RequireUserAdmin(u.EffectiveChat, u.EffectiveMessage, u.EffectiveUser.Id, nil) {
		return gotgbot.ContinueGroups{}
	}

	if len(args) >= 1 {
		noteName := args[0]

		if sql.RmNote(strconv.Itoa(chatId), noteName) {
			_, err := u.EffectiveMessage.ReplyTextf("Successfully removed note %v", noteName)
			return err
		} else {
			_, err := u.EffectiveMessage.ReplyText("That's not a note in my database!")
			return err
		}
	}
	return nil
}

func listNotes(_ ext.Bot, u *gotgbot.Update) error {
	chatId := u.EffectiveChat.Id
	noteList := sql.GetAllChatNotes(strconv.Itoa(chatId))

	msg := "<code>Notes in chat:</code>\n"
	for _, note := range noteList {
		noteName := html.EscapeString(fmt.Sprintf(" - %v\n", note.Name))
		if len(msg) + len(noteName) > helpers.MaxMessageLength {
			_, err := u.EffectiveMessage.ReplyHTML(msg)
			msg = ""
			error_handling.HandleErr(err)
		}
		msg += noteName
	}

	if msg == "<code>Notes in chat:</code>\n" {
		_, err := u.EffectiveMessage.ReplyText("No notes in this chat!")
		return err
	} else if len(msg) != 0 {
		_, err := u.EffectiveMessage.ReplyHTML(msg)
		return err
	}
	return nil
}

func LoadNotes(u *gotgbot.Updater) {
	defer log.Println("Loading module notes")
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("get", cmdGet))
	u.Dispatcher.AddHandler(handlers.NewRegex(`^#[^\s]+`, hashGet))
	u.Dispatcher.AddHandler(handlers.NewCommand("save", save))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("clear", clear))
	u.Dispatcher.AddHandler(handlers.NewCommand("notes", listNotes))
	u.Dispatcher.AddHandler(handlers.NewCommand("saved", listNotes))
}
