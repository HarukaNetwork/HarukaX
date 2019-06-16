package users

import (
	"github.com/ATechnoHazard/ginko/go_bot/modules/sql"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/handlers/Filters"
	"log"
	"strconv"
)

func logUsers(_ ext.Bot, u *gotgbot.Update) error {
	chat := u.EffectiveChat
	msg := u.EffectiveMessage

	sql.UpdateUser(msg.From.Id,
		msg.From.Username,
		strconv.Itoa(chat.Id),
		chat.Title)

	if msg.ReplyToMessage != nil {
		sql.UpdateUser(msg.From.Id,
			msg.From.Username,
			strconv.Itoa(chat.Id),
			chat.Title)
	}

	if msg.ForwardFrom != nil {
		sql.UpdateUser(msg.ForwardFrom.Id,
			msg.ForwardFrom.Username, "nil", "nil")
	}

	return gotgbot.ContinueGroups{}
}

func GetUserId(username string) int {
	if len(username) <= 5 {
		return 0
	}
	if username[0] == '@' {
		username = username[1:]
	}
	users := sql.GetUserIdByName(username)
	if users == nil {
		return 0
	}

	return users.UserId
}

func LoadUsers(u *gotgbot.Updater) {
	defer log.Println("Loading module users")
	u.Dispatcher.AddHandler(handlers.NewMessage(Filters.All, logUsers))
}
