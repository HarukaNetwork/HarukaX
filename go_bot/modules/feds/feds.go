package feds

import (
	"github.com/ATechnoHazard/ginko/go_bot/modules/sql"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/google/uuid"
	"log"
	"strconv"
	"strings"
)

func NewFed(_ ext.Bot, u *gotgbot.Update) error {
	user := u.EffectiveUser
	msg := u.EffectiveMessage

	splitText := strings.SplitN(msg.Text, " ", 2)
	if len(splitText) < 2 {
		_, err := msg.ReplyText("Please send me the name of the federation you want to create!")
		return err
	}

	fedName := splitText[1]

	fedId := uuid.New().String()

	fed := sql.NewFed(strconv.Itoa(user.Id), fedId, fedName)
	if fed == "" {
		_, err := msg.ReplyText("Big F! Couldn't create a new federation.")
		return err
	}
	_, err := msg.ReplyHTMLf("<b>You have successfully created a new federation!</b>"+
		"\nName: <code>%v</code>"+
		"\nID: <code>%v</code>"+
		"\nUse the following command to join the federation:"+
		"\n<code>/joinfed %v</code>", fedName, fed, fed)
	return err
}

func DelFed(_ ext.Bot, u *gotgbot.Update) error {
	user := u.EffectiveUser
	msg := u.EffectiveMessage

	if u.EffectiveChat.Type != "private" {
		_, err := msg.ReplyText("Delete your federation in my PM - not in a group.")
		return err
	}

	fed:= sql.GetFedFromUser(strconv.Itoa(user.Id))

	if fed == nil {
		_, err := msg.ReplyText("You aren't the creator of any federations!")
		return err
	}

	go sql.DelFed(fed.FedId)
	_, err := msg.ReplyHTMLf("Federation <b>%v</b> has been deleted!", fed.FedName)
	return err
}

func LoadFeds(u *gotgbot.Updater) {
	defer log.Println("Loading module feds")
	u.Dispatcher.AddHandler(handlers.NewCommand("newfed", NewFed))
	u.Dispatcher.AddHandler(handlers.NewCommand("delfed", DelFed))
}
