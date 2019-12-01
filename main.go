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

package main

import (
	"github.com/ATechnoHazard/ginko/go_bot"
	"github.com/ATechnoHazard/ginko/go_bot/modules/admin"
	"github.com/ATechnoHazard/ginko/go_bot/modules/bans"
	"github.com/ATechnoHazard/ginko/go_bot/modules/blacklist"
	"github.com/ATechnoHazard/ginko/go_bot/modules/deleting"
	"github.com/ATechnoHazard/ginko/go_bot/modules/feds"
	"github.com/ATechnoHazard/ginko/go_bot/modules/help"
	"github.com/ATechnoHazard/ginko/go_bot/modules/misc"
	"github.com/ATechnoHazard/ginko/go_bot/modules/muting"
	"github.com/ATechnoHazard/ginko/go_bot/modules/notes"
	"github.com/ATechnoHazard/ginko/go_bot/modules/sql"
	"github.com/ATechnoHazard/ginko/go_bot/modules/users"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/error_handling"
	"github.com/ATechnoHazard/ginko/go_bot/modules/warns"
	"github.com/ATechnoHazard/ginko/go_bot/modules/welcome"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"log"
)

func main() {
	// Create updater instance
	u, err := gotgbot.NewUpdater(go_bot.BotConfig.ApiKey)
	error_handling.FatalError(err)

	// Add start handler
	u.Dispatcher.AddHandler(handlers.NewCommand("start", start))

	// Create database tables if not already existing
	sql.EnsureBotInDb(u)

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

	log.Println("Starting long polling")
	err = u.StartPolling()
	error_handling.HandleErr(err)
	u.Idle()
}

func start(_ ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage
	_, err := msg.ReplyTextf("Hi there! I'm a telegram group management bot, written in Go." +
		"\nFor any questions or bug reports, you can head over to @gobotsupport.")
	error_handling.HandleErr(err)
	return nil
}
