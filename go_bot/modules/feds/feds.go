package feds

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"log"
)

func LoadFeds(u *gotgbot.Updater) {
	defer log.Println("Loading module feds")
	u.Dispatcher.AddHandler(handlers.NewCommand("newfed", newFed))
	u.Dispatcher.AddHandler(handlers.NewCommand("delfed", delFed))
	u.Dispatcher.AddHandler(handlers.NewCommand("chatfed", chatFed))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("joinfed", joinFed))
	u.Dispatcher.AddHandler(handlers.NewCommand("leavefed", leaveFed))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("fedpromote", fedPromote))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("feddemote", fedDemote))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("fedinfo", fedInfo))
	u.Dispatcher.AddHandler(handlers.NewArgsCommand("fedadmins", fedAdmins))
}
