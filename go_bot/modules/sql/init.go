package sql

import (
	"github.com/ATechnoHazard/ginko/go_bot"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/error_handling"
	"github.com/go-pg/pg"
)

var SESSION *pg.DB

func init() {
	opt, err := pg.ParseURL(go_bot.BotConfig.SqlUri)
	error_handling.HandleErrorAndExit(err)
	if go_bot.BotConfig.Heroku {
		opt.PoolSize = 20
	}
	SESSION = pg.Connect(opt)
}
