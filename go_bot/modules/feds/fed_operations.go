/*
 *   Copyright 2019 ATechnoHazard  <amolele@gmail.com>
 *
 *   Permission is hereby granted, free of charge, to any person obtaining a copy
 *   of this software and associated documentation files (the "Software"), to deal
 *   in the Software without restriction, including without limitation the rights
 *   to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *   copies of the Software, and to permit persons to whom the Software is
 *   furnished to do so, subject to the following conditions:
 *
 *   The above copyright notice and this permission notice shall be included in all
 *   copies or substantial portions of the Software.
 *
 *   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *   IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *   FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *   AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *   LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *   OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *   SOFTWARE.
 */

package feds

import (
	"github.com/ATechnoHazard/ginko/go_bot/modules/sql"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/error_handling"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/extraction"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/helpers"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/google/uuid"
	"strconv"
	"strings"
)

func newFed(_ ext.Bot, u *gotgbot.Update) error {
	user := u.EffectiveUser
	msg := u.EffectiveMessage
	var fedId string

	splitText := strings.SplitN(msg.Text, " ", 2)
	if len(splitText) < 2 {
		_, err := msg.ReplyText("Please send me the name of the federation you want to create!")
		return err
	}

	fedName := splitText[1]

	existingFed := sql.GetFedFromOwnerId(strconv.Itoa(user.Id))

	if existingFed != nil {
		fedId = existingFed.Id
	} else {
		fedId = uuid.New().String()
	}

	fed := sql.NewFed(strconv.Itoa(user.Id), fedId, fedName)
	if !fed {
		_, err := msg.ReplyText("Big F! Couldn't create a new federation.")
		return err
	}
	_, err := msg.ReplyHTMLf("<b>You have successfully created a new federation!</b>"+
		"\nName: <code>%v</code>"+
		"\nID: <code>%v</code>"+
		"\nUse the following command to join the federation:"+
		"\n<code>/joinfed %v</code>", fedName, fedId, fedId)
	return err
}

func delFed(_ ext.Bot, u *gotgbot.Update) error {
	user := u.EffectiveUser
	msg := u.EffectiveMessage

	if u.EffectiveChat.Type != "private" {
		_, err := msg.ReplyText("Delete your federation in my PM - not in a group.")
		return err
	}

	fed := sql.GetFedFromOwnerId(strconv.Itoa(user.Id))

	if fed == nil {
		_, err := msg.ReplyText("You aren't the creator of any federations!")
		return err
	}

	go sql.DelFed(fed.Id)
	_, err := msg.ReplyHTMLf("Federation <b>%v</b> has been deleted!", fed.FedName)
	return err
}

func joinFed(_ ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat
	user := u.EffectiveUser
	msg := u.EffectiveMessage

	member, err := chat.GetMember(user.Id)
	error_handling.HandleErr(err)

	if member.Status != "creator" {
		_, err := msg.ReplyText("Only group creators can join federations.")
		return err
	}

	fedId := sql.GetFedId(strconv.Itoa(chat.Id))
	if fedId != "" {
		_, err := msg.ReplyText("You cannot join more that one federation per chat!")
		return err
	}

	if len(args) >= 1 {
		toJoinId := args[0]
		fed := sql.GetFedInfo(toJoinId)
		if fed == nil {
			_, err := msg.ReplyText("Please enter a valid federation ID!")
			return err
		}

		result := sql.ChatJoinFed(toJoinId, strconv.Itoa(chat.Id))
		if !result {
			_, err := msg.ReplyText("Well damn, I couldn't join that fed.")
			return err
		}
		_, err := msg.ReplyHTMLf("Joined federation <b>%v</b>.", fed.FedName)
		return err
	} else {
		_, err := msg.ReplyText("Please enter a federation ID!")
		return err
	}
}

func leaveFed(_ ext.Bot, u *gotgbot.Update) error {
	chat := u.EffectiveChat
	user := u.EffectiveUser
	msg := u.EffectiveMessage

	member, err := chat.GetMember(user.Id)
	error_handling.HandleErr(err)

	if member.Status != "creator" {
		_, err := msg.ReplyText("Only group creators can leave federations.")
		return err
	}

	if sql.ChatLeaveFed(strconv.Itoa(chat.Id)) {
		_, err := msg.ReplyHTMLf("Left federation!")
		return err
	} else {
		_, err := msg.ReplyHTMLf("This chat isn't in any federations!")
		return err
	}
}

func fedPromote(bot ext.Bot, u *gotgbot.Update, args []string) error {
	user := u.EffectiveUser
	msg := u.EffectiveMessage

	uId := extraction.ExtractUser(msg, args)
	userId := strconv.Itoa(uId)
	if userId == "0" {
		_, err := msg.ReplyText("Try targeting a user next time bud.")
		return err
	}

	member, err := bot.GetChat(uId)
	error_handling.HandleErr(err)

	fed := sql.GetFedFromOwnerId(strconv.Itoa(user.Id))
	if fed == nil {
		_, err := msg.ReplyText("You aren't the creator of any federations.")
		return err
	}

	if userId == fed.OwnerId {
		_, err := msg.ReplyText("Why are you trying to promote yourself in your own federation?")
		return err
	}

	if sql.IsUserFedAdmin(fed.Id, userId) != "" {
		_, err := msg.ReplyText("This user is already a federation admin in your federation.")
		return err
	}

	if userId == strconv.Itoa(bot.Id) {
		_, err := msg.ReplyText("I don't need to be an admin in any feds!")
		return err
	}

	go sql.UserPromoteFed(fed.Id, userId)

	_, err = msg.ReplyHTMLf("User %v is now an admin of <b>%v</b> (<code>%v</code>)", helpers.MentionHtml(uId, member.FirstName), fed.FedName, fed.Id)
	return err
}

func fedDemote(bot ext.Bot, u *gotgbot.Update, args []string) error {
	user := u.EffectiveUser
	msg := u.EffectiveMessage

	uId := extraction.ExtractUser(msg, args)
	userId := strconv.Itoa(uId)
	if userId == "0" {
		_, err := msg.ReplyText("Try targeting a user next time bud.")
		return err
	}

	member, err := bot.GetChat(uId)
	error_handling.HandleErr(err)

	fed := sql.GetFedFromOwnerId(strconv.Itoa(user.Id))
	if fed == nil {
		_, err := msg.ReplyText("You aren't the creator of any federations.")
		return err
	}

	if userId == fed.OwnerId {
		_, err := msg.ReplyText("Why are you trying to demote yourself in your own federation?")
		return err
	}

	if sql.IsUserFedAdmin(fed.Id, userId) == "" {
		_, err := msg.ReplyText("This user is not a federation admin in your federation.")
		return err
	}

	go sql.UserDemoteFed(fed.Id, userId)

	_, err = msg.ReplyHTMLf("User %v is no longer an admin of <b>%v</b> (<code>%v</code>)", helpers.MentionHtml(uId, member.FirstName), fed.FedName, fed.Id)
	return err
}
