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

package main

import (
	"log"
	"strconv"

	"github.com/HarukaNetwork/HarukaX/harukax/modules/rules"

	"github.com/HarukaNetwork/HarukaX/harukax"
	"github.com/HarukaNetwork/HarukaX/harukax/modules/admin"
	"github.com/HarukaNetwork/HarukaX/harukax/modules/bans"
	"github.com/HarukaNetwork/HarukaX/harukax/modules/blacklist"
	"github.com/HarukaNetwork/HarukaX/harukax/modules/deleting"
	"github.com/HarukaNetwork/HarukaX/harukax/modules/feds"
	"github.com/HarukaNetwork/HarukaX/harukax/modules/help"
	"github.com/HarukaNetwork/HarukaX/harukax/modules/misc"
	"github.com/HarukaNetwork/HarukaX/harukax/modules/muting"
	"github.com/HarukaNetwork/HarukaX/harukax/modules/notes"
	"github.com/HarukaNetwork/HarukaX/harukax/modules/sql"
	"github.com/HarukaNetwork/HarukaX/harukax/modules/users"
	"github.com/HarukaNetwork/HarukaX/harukax/modules/utils/caching"
	"github.com/HarukaNetwork/HarukaX/harukax/modules/utils/error_handling"
	"github.com/HarukaNetwork/HarukaX/harukax/modules/warns"
	"github.com/HarukaNetwork/HarukaX/harukax/modules/welcome"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
)

func main() {
	// Create updater instance
	u, err := gotgbot.NewUpdater(harukax.BotConfig.ApiKey)
	error_handling.FatalError(err)

	// Add start handler
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("start", start))

	// Create database tables if not already existing
	sql.EnsureBotInDb(u)

	// Prepare Caching Service
	caching.InitCache()
	//caching.InitRedis()

	// Add module handlers
	bans.LoadBans(u)
	users.LoadUsers(u)
	admin.LoadAdmin(u)
	warns.LoadWarns(u)
	misc.LoadMisc(u)
	muting.LoadMuting(u)
	deleting.LoadDelete(u)
	blacklist.LoadBlacklist(u)
	feds.LoadFeds(u)
	notes.LoadNotes(u)
	help.LoadHelp(u)
	welcome.LoadWelcome(u)
	rules.LoadRules(u)

	if harukax.BotConfig.DropUpdate == "True" {
		log.Println("[Info][Core] Using Clean Long Polling")
		err = u.StartCleanPolling()
		error_handling.HandleErr(err)
	} else {
		log.Println("[Info][Core] Using Long Polling")
		err = u.StartPolling()
		error_handling.HandleErr(err)
	}

	u.Idle()
}

func start(_ ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage

	if u.EffectiveChat.Type == "private" {
		if len(args) != 0 {
			if _, err := strconv.Atoi(args[0][2:]); err == nil {
				chatRules := sql.GetChatRules(args[0])
				if chatRules != nil {
					_, err := msg.ReplyHTML(chatRules.Rules)
					return err
				}
				_, err := msg.ReplyText("The group admins haven't set any rules for this chat yet. This probably doesn't " +
					"mean it's lawless though!")
				log.Println(args[0])
				return err
			}
		}
	}

	_, err := msg.ReplyTextf("Hi there! I'm a telegram group management bot, written in Go." +
		"\nFor any questions or bug reports, you can head over to @gobotsupport.")
	return err
}
