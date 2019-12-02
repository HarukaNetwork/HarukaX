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

package sql

import (
	"log"

	"github.com/ATechnoHazard/ginko/go_bot"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/error_handling"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
)

var SESSION *gorm.DB

func init() {
	conn, err := pq.ParseURL(go_bot.BotConfig.SqlUri)
	error_handling.FatalError(err)

	db, err := gorm.Open("postgres", conn)
	error_handling.FatalError(err)

	if go_bot.BotConfig.DebugMode {
		SESSION = db.Debug()
		log.Println("Using database in debug mode.")
	} else {
		SESSION = db
	}

	if go_bot.BotConfig.Heroku {
		db.DB().SetMaxOpenConns(20)
	} else {
		db.DB().SetMaxOpenConns(100)
	}

	log.Println("Database connected")

	// Create tables if they don't exist
	SESSION.AutoMigrate(&User{}, &Chat{}, &Warns{}, &WarnFilters{}, &WarnSettings{}, &BlackListFilters{}, &Federation{},
		&FedChat{}, &FedAdmin{}, &FedBan{}, &Note{}, &Button{}, &Welcome{}, &WelcomeButton{}, &MutedUser{})
	log.Println("Auto-migrated database schema")

}
