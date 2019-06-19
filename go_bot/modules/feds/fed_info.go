package feds

import (
	"fmt"
	"github.com/ATechnoHazard/ginko/go_bot/modules/sql"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/helpers"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"strconv"
)

func chatFed(_ ext.Bot, u *gotgbot.Update) error {
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	if chat.Type == "private" {
		_, err := msg.ReplyText("Your PM cannot be part of a federation!")
		return err
	}

	fedId := sql.GetFedId(strconv.Itoa(chat.Id))

	if fedId == "" {
		_, err := msg.ReplyText("This chat is not part of a federation.")
		return err
	}

	fed := sql.GetFedInfo(fedId)
	_, err := msg.ReplyHTMLf("This chat is part of the following federation:" +
		"\n<b>%v</b> (ID: <code>%v</code>)", fed.FedName, fedId)
	return err
}

func fedInfo(_ ext.Bot, u *gotgbot.Update, args []string) error {
	user := u.EffectiveUser
	msg := u.EffectiveMessage
	var fedId string
	var fed *sql.Federations
	if len(args) >= 1 {
		fedId = args[0]
		fed = sql.GetFedInfo(fedId)
		if fed == nil {
			_, err := msg.ReplyText("Please use a valid federation ID.")
			return err
		}
	} else {
		fed = sql.GetFedFromUser(strconv.Itoa(user.Id))
		if fed == nil {
			_, err := msg.ReplyText("You aren't the creator of any federations!")
			return err
		}
		fedId = fed.FedId

	}

	ownerId, _ := strconv.Atoi(fed.OwnerId)

	text := fmt.Sprintf("<b>Fed info:</b>" +
		"\nFedID: <code>%v</code>" +
		"\nName: %v" +
		"\nCreator: %v" +
		"\nNumber of admins: <code>%v</code>" +
		"\nNumber of bans: <code>%v</code>" +
		"\nNumber of connected chats: <code>%v</code>", fed.FedId, fed.FedName, helpers.MentionHtml(ownerId, "this person"),
		len(fed.FedAdmins) + 1,
		len(sql.GetAllFbanUsers(fedId)),
		len(sql.AllFedChats(fedId)))

	_, err := msg.ReplyHTML(text)
	return err
}

func fedAdmins(bot ext.Bot, u *gotgbot.Update, args []string) error {
	user := u.EffectiveUser
	msg := u.EffectiveMessage
	var fedId string
	var fed *sql.Federations
	if len(args) >= 1 {
		fedId = args[0]
		fed = sql.GetFedInfo(fedId)
		if fed == nil {
			_, err := msg.ReplyText("Please use a valid federation ID.")
			return err
		}
	} else {
		fed = sql.GetFedFromUser(strconv.Itoa(user.Id))
		fedId = fed.FedId
		if fed == nil {
			_, err := msg.ReplyText("You aren't the creator of any federations!")
			return err
		}
	}

	ownerId, _ := strconv.Atoi(fed.OwnerId)
	owner, _ := bot.GetChat(ownerId)

	text := "Admins in this federation:"
	text += fmt.Sprintf("\n - %v (<code>%v</code>)", helpers.MentionHtml(ownerId, owner.FirstName), ownerId)
	for _, users := range fed.FedAdmins {
		userId, _ := strconv.Atoi(users)
		user, _ := bot.GetChat(userId)
		text += fmt.Sprintf("\n - %v (<code>%v</code>)", helpers.MentionHtml(user.Id, user.FirstName), user.Id)
	}

	_, err := msg.ReplyHTML(text)
	return err
}
