package sql

import (
	"github.com/ATechnoHazard/ginko/go_bot"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/error_handling"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
	"log"
)

var SESSION *gorm.DB

func init() {
	conn, err := pq.ParseURL(go_bot.BotConfig.SqlUri)
	error_handling.FatalError(err)

	db, err := gorm.Open("postgres", conn)
	error_handling.FatalError(err)
	SESSION = db

	if go_bot.BotConfig.Heroku {
		db.DB().SetMaxOpenConns(20)
	} else {
		db.DB().SetMaxOpenConns(100)
	}

	log.Println("Database connected")

	// Create tables if they don't exist
	SESSION.AutoMigrate(&User{}, &Chat{}, &Warns{}, &WarnFilters{}, &WarnSettings{}, &BlackListFilters{}, &Federation{},
	&FedChat{}, &FedAdmin{}, &FedBan{}, &Note{}, &Button{})
	log.Println("Auto-migrated database schema")
}
