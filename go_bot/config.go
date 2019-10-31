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

package go_bot

import (
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/error_handling"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	BotName     string
	ApiKey      string
	OwnerName   string
	OwnerId     int
	SudoUsers   []string
	LoadPlugins []string
	SqlUri      string
	Heroku      bool
}

var BotConfig Config

// Returns a config object generated from the dotenv file
func init() {
	err := godotenv.Load()
	error_handling.FatalError(err)
	returnConfig := Config{}

	// Assign config struct values by loading them from the env
	var ok bool
	returnConfig.BotName, ok = os.LookupEnv("BOT_NAME")
	// If env var is empty
	if !ok {
		log.Fatal("Missing bot name")
	}

	returnConfig.ApiKey, ok = os.LookupEnv("BOT_API_KEY") // If env var is empty
	if !ok {
		log.Fatal("Missing API key")
	}

	returnConfig.OwnerName, ok = os.LookupEnv("OWNER_USERNAME")
	// If env var is empty
	if !ok {
		log.Fatal("Missing owner username")
	}

	returnConfig.OwnerId, err = strconv.Atoi(os.Getenv("OWNER_ID"))
	error_handling.FatalError(err)

	returnConfig.SudoUsers = strings.Split(os.Getenv("SUDO_USERS"), " ")

	returnConfig.SqlUri, ok = os.LookupEnv("DATABASE_URI")
	// If env var is empty
	if !ok {
		log.Fatal("Missing PostgreSQL URI")
	}

	_, returnConfig.Heroku = os.LookupEnv("HEROKU")

	BotConfig = returnConfig

}
