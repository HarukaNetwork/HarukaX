package main

import (
	"github.com/ATechnoHazard/ginko/go_bot"
	"github.com/ATechnoHazard/ginko/go_bot/modules/admin"
	"github.com/ATechnoHazard/ginko/go_bot/modules/bans"
	"github.com/ATechnoHazard/ginko/go_bot/modules/misc"
	"github.com/ATechnoHazard/ginko/go_bot/modules/muting"
	"github.com/ATechnoHazard/ginko/go_bot/modules/sql"
	"github.com/ATechnoHazard/ginko/go_bot/modules/users"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/error_handling"
	"github.com/ATechnoHazard/ginko/go_bot/modules/warns"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"log"
)

func main() {
	log.Println("Starting long polling")
	u, err := gotgbot.NewUpdater(go_bot.BotConfig.ApiKey)
	error_handling.HandleErrorAndExit(err)

	// Add module handlers
	bans.LoadBans(u)
	users.LoadUsers(u)
	admin.LoadAdmin(u)
	warns.LoadWarns(u)
	misc.LoadMisc(u)
	muting.LoadMuting(u)

	sql.EnsureBotInDb(u)

	u.Dispatcher.AddHandler(handlers.NewCommand("start", start))

	err = u.StartPolling()
	error_handling.HandleErr(err)
	u.Idle()
}

func start(_ ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage
	_, err := msg.ReplyText("Hi there!")
	error_handling.HandleErr(err)
	return nil
}
