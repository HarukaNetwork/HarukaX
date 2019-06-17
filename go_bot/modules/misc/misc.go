package misc

import (
	"fmt"
	"github.com/ATechnoHazard/ginko/go_bot"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/error_handling"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/extraction"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/helpers"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"html"
	"log"
	"strconv"
	"strings"
)

func getId(bot ext.Bot, u *gotgbot.Update, args []string) error {
	userId := extraction.ExtractUser(u.EffectiveMessage, args)
	if userId != 0 {
		if u.EffectiveMessage.ReplyToMessage != nil && u.EffectiveMessage.ReplyToMessage.ForwardFrom != nil {
			user1 := u.EffectiveMessage.ReplyToMessage.From
			user2 := u.EffectiveMessage.ReplyToMessage.ForwardFrom
			_, err := u.EffectiveMessage.ReplyHTMLf("The original sender, %v, has an ID of <code>%v</code>.\n"+
				"The forwarder, %v, has an ID of <code>%v</code>.", html.EscapeString(user2.FirstName),
				user2.Id,
				html.EscapeString(user1.FirstName),
				user1.Id)
			return err
		} else {
			user, err := bot.GetChat(userId)
			error_handling.HandleErr(err)
			_, err = u.EffectiveMessage.ReplyHTMLf("%v's ID is <code>%v</code>", html.EscapeString(user.FirstName), user.Id)
		}
	} else {
		chat := u.EffectiveChat
		if chat.Type == "private" {
			_, err := u.EffectiveMessage.ReplyHTMLf("Your ID is <code>%v</code>", chat.Id)
			return err
		} else {
			_, err := u.EffectiveMessage.ReplyHTMLf("This group's ID is <code>%v</code>", chat.Id)
			return err
		}
	}
	return nil
}

func info(bot ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	userId := extraction.ExtractUser(msg, args)
	var user *ext.User

	if userId != 0 {
		userChat, _ := bot.GetChat(userId)
		user = &ext.User{
			Id:        userChat.Id,
			FirstName: userChat.FirstName,
			LastName:  userChat.LastName,
		}

	} else if msg.ReplyToMessage == nil && len(args) <= 0 {
		user = msg.From
		userId = msg.From.Id

	} else if _, err := strconv.Atoi(args[0]); msg.ReplyToMessage == nil && (len(args) <= 0 || (len(args) >= 1 && strings.HasPrefix(args[0], "@") && err != nil && msg.ParseEntities()[0].Type != "TEXT_MENTION")) {
		_, err := msg.ReplyText("Yeah nah, this mans doesn't exist.")
		return err
	} else {
		return nil
	}

	text := fmt.Sprintf("<b>User info</b>"+
		"\nID: <code>%v</code>"+
		"\nFirst Name: %v", userId, html.EscapeString(user.FirstName))

	if user.LastName != "" {
		text += fmt.Sprintf("\nLast Name: %v", user.LastName)
	}

	if user.Username != "" {
		text += fmt.Sprintf("\nUsername: @%v", user.Username)
	}

	text += fmt.Sprintf("\nPermanent user link: %v", helpers.MentionHtml(user.Id, user.FirstName+user.LastName))

	if user.Id == go_bot.BotConfig.OwnerId {
		text += "\n\nDis nibba stronk af!"
	} else {
		for _, id := range go_bot.BotConfig.SudoUsers {
			if strconv.Itoa(user.Id) == id {
				text += "\n\nThis person is one of my sudo users! " +
					"Nearly as powerful as my owner - so watch it."
			}
		}
	}
	_, err := u.EffectiveMessage.ReplyHTML(text)
	return err
}

func LoadMisc(u *gotgbot.Updater) {
	defer log.Println("Loading module misc")
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("id", getId))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("info", info))
}
