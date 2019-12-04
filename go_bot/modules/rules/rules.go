package rules

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ATechnoHazard/ginko/go_bot/modules/sql"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/helpers"
	tg_md2html "github.com/PaulSonOfLars/gotg_md2html"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
)

func sendRules(bot ext.Bot, u *gotgbot.Update) error {
	rules := sql.GetChatRules(strconv.Itoa(u.EffectiveChat.Id))
	log.Println(rules)

	if rules != nil {
		if rules.Rules != "" {
			msg := bot.NewSendableMessage(u.EffectiveChat.Id, "Contact me in PM to get this group's rules.")
			button := sql.WelcomeButton{
				Id:       0,
				ChatId:   strconv.Itoa(u.EffectiveChat.Id),
				Name:     "Rules",
				Url:      fmt.Sprintf("t.me/%v?start=%v", bot.UserName, u.EffectiveChat.Id),
				SameLine: false,
			}
			keyb := helpers.BuildWelcomeKeyboard([]sql.WelcomeButton{button})
			keyboard := ext.InlineKeyboardMarkup{InlineKeyboard: &keyb}
			msg.ReplyMarkup = &keyboard
			_, err := msg.Send()
			return err
		}
	}
	_, err := u.EffectiveMessage.ReplyText("The group admins haven't set any rules for this chat yet. This probably doesn't " +
		"mean it's lawless though!")
	return err
}

func setRules(_ ext.Bot, u *gotgbot.Update) error {
	chatId := strconv.Itoa(u.EffectiveChat.Id)
	msg := u.EffectiveMessage
	rawText := msg.Text
	timesInserted := 0
	entities := msg.Entities

	for _, ent := range entities {
		if ent.Type == "code" {
			rawText = rawText[:ent.Offset+timesInserted] + "`" + rawText[ent.Offset+timesInserted:]
			timesInserted++
			rawText = rawText[:(ent.Offset+ent.Length+(timesInserted))] + "`" + rawText[(ent.Offset+ent.Length+(timesInserted)):]
			timesInserted++
		}
	}
	args := strings.SplitN(rawText, " ", 2)
	if len(args) == 2 {
		txt := tg_md2html.MD2HTML(args[1])
		go sql.SetChatRules(chatId, txt)
		_, err := msg.ReplyText("Successfully set rules for this group!")
		return err
	}

	_, err := msg.ReplyText("You need to give me some rules to set!")
	return err
}

func clearRules(_ ext.Bot, u *gotgbot.Update) error {
	chatId := strconv.Itoa(u.EffectiveChat.Id)
	go sql.SetChatRules(chatId, "")
	_, err := u.EffectiveMessage.ReplyText("Successfully cleared rules!")
	return err
}

func LoadRules(u *gotgbot.Updater) {
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("rules", []rune{'/', '!'}, sendRules))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("setrules", []rune{'/', '!'}, setRules))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("clearrules", []rune{'/', '!'}, clearRules))
}
