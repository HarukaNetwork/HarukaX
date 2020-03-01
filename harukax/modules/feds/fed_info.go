/*
 *    Copyright © 2020 Haruka Network Development
 *    This file is part of Haruka X.
 *
 *    Haruka X is free software: you can redistribute it and/or modify
 *    it under the terms of the Raphielscape Public License as published by
 *    the Devscapes Open Source Holding GmbH., version 1.d
 *
 *    Haruka X is distributed in the hope that it will be useful,
 *    but WITHOUT ANY WARRANTY; without even the implied warranty of
 *    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *    Devscapes Raphielscape Public License for more details.
 *
 *    You should have received a copy of the Devscapes Raphielscape Public License
 */

package feds

import (
	"fmt"
	"strconv"

	"github.com/HarukaNetwork/HarukaX/harukax/modules/sql"
	"github.com/HarukaNetwork/HarukaX/harukax/modules/utils/extraction"
	"github.com/HarukaNetwork/HarukaX/harukax/modules/utils/helpers"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
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
	_, err := msg.ReplyHTMLf("This chat is part of the following federation:"+
		"\n<b>%v</b> (ID: <code>%v</code>)", fed.FedName, fedId)
	return err
}

func fedInfo(_ ext.Bot, u *gotgbot.Update, args []string) error {
	user := u.EffectiveUser
	msg := u.EffectiveMessage
	var fedId string
	var fed *sql.Federation
	if len(args) >= 1 {
		fedId = args[0]
		fed = sql.GetFedInfo(fedId)
		if fed == nil {
			_, err := msg.ReplyText("Please use a valid federation ID.")
			return err
		}
	} else {
		fed = sql.GetFedFromOwnerId(strconv.Itoa(user.Id))
		if fed == nil {
			_, err := msg.ReplyText("You aren't the creator of any federations!")
			return err
		}
		fedId = fed.Id
	}

	ownerId, _ := strconv.Atoi(fed.OwnerId)

	text := fmt.Sprintf("<b>Fed info:</b>"+
		"\nFedID: <code>%v</code>"+
		"\nName: %v"+
		"\nCreator: %v"+
		"\nNumber of admins: <code>%v</code>"+
		"\nNumber of bans: <code>%v</code>"+
		"\nNumber of connected chats: <code>%v</code>", fed.Id, fed.FedName, helpers.MentionHtml(ownerId, "this person"),
		len(sql.GetFedAdmins(fedId))+1,
		sql.GetFbanUsersCount(fedId),
		len(sql.AllFedChats(fedId)))

	_, err := msg.ReplyHTML(text)
	return err
}

func fedAdmins(bot ext.Bot, u *gotgbot.Update, args []string) error {
	user := u.EffectiveUser
	msg := u.EffectiveMessage
	var fedId string
	var fed *sql.Federation
	if len(args) >= 1 {
		fedId = args[0]
		fed = sql.GetFedInfo(fedId)
		if fed == nil {
			_, err := msg.ReplyText("Please use a valid federation ID.")
			return err
		}
	} else {
		fed = sql.GetFedFromOwnerId(strconv.Itoa(user.Id))
		fedId = fed.Id
		if fed == nil {
			_, err := msg.ReplyText("You aren't the creator of any federations!")
			return err
		}
	}

	ownerId, _ := strconv.Atoi(fed.OwnerId)
	owner, _ := bot.GetChat(ownerId)

	text := "Admins in this federation:"
	text += fmt.Sprintf("\n - %v (<code>%v</code>)", helpers.MentionHtml(ownerId, owner.FirstName), ownerId)
	for _, users := range sql.GetFedAdmins(fedId) {
		userId, _ := strconv.Atoi(users.UserId)
		user, _ := bot.GetChat(userId)
		text += fmt.Sprintf("\n - %v (<code>%v</code>)", helpers.MentionHtml(user.Id, user.FirstName), user.Id)
	}

	_, err := msg.ReplyHTML(text)
	return err
}

func fedStat(bot ext.Bot, u *gotgbot.Update, args []string) error {
	user := u.EffectiveUser
	msg := u.EffectiveMessage
	var fedId = ""

	userId := extraction.ExtractUser(msg, args)
	if userId == 0 {
		userId = user.Id
		if len(args) >= 1 {
			fedId = args[0]
		}
	} else {
		if len(args) >= 2 {
			fedId = args[1]
		}
	}

	chatMember, err := bot.GetChat(userId)
	if err != nil {
		return err
	}

	if fedId != "" {
		fed := sql.GetFedInfo(fedId)
		if fed == nil {
			_, err := msg.ReplyText("Please use a valid federation ID!")
			return err
		}

		fban := sql.GetFbanUser(fedId, strconv.Itoa(userId))
		if fban == nil {
			_, err := msg.ReplyText("Good news! You aren't fedbanned in this federation!")
			return err
		} else {
			text := fmt.Sprintf("%v is fedbanned in <b>%v</b> for the following reason:\n"+
				" - %v", helpers.MentionHtml(chatMember.Id, chatMember.FirstName), fed.FedName, fban.Reason)
			_, err := msg.ReplyHTML(text)
			return err
		}
	}

	feds := sql.GetUserFbans(strconv.Itoa(userId))
	if feds == nil {
		_, err := msg.ReplyText("Well damn, something went wrong.")
		return err
	} else if len(feds) == 0 {
		_, err := msg.ReplyHTMLf("%v hasn't been fedbanned anywhere!", helpers.MentionHtml(userId, chatMember.FirstName))
		return err
	}
	text := fmt.Sprintf("%v has been banned in the following federations:\n", helpers.MentionHtml(userId, chatMember.FirstName))
	for _, fed := range feds {
		text += fmt.Sprintf("<b>%v</b> - <code>%v</code>\n", fed.FedName, fed.Id)
	}
	text += "If you would like to know more about the fedban reason in a specific federation, use /fbanstat FedID."

	_, err = msg.ReplyHTMLf(text)
	return err
}
