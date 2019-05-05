package main

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/atechnohazard/ginko/go_bot"
	"github.com/atechnohazard/ginko/go_bot/modules/bans"
	"github.com/atechnohazard/ginko/go_bot/modules/sql"
	"github.com/atechnohazard/ginko/go_bot/modules/users"
	"github.com/atechnohazard/ginko/go_bot/modules/utils/error_handling"
	"log"
)

func main() {
	log.Println("Starting long polling")
	u, err := gotgbot.NewUpdater(go_bot.BotConfig.ApiKey)
	error_handling.HandleErrorAndExit(err)

	// Add module handlers
	bans.LoadBans(u)
	users.LoadUsers(u)

	sql.EnsureBotInDb(u)

	u.Dispatcher.AddHandler(handlers.NewCommand("start", start))

	err = u.StartPolling()
	error_handling.HandleErrorGracefully(err)
	u.Idle()
}

func start(_ ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage
	_, err := msg.ReplyText("Hewwo, chu started me desu~")
	error_handling.HandleErrorGracefully(err)
	return nil
}
