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
	BotName string
	ApiKey string
	OwnerName string
	OwnerId int
	SudoUsers []string
	LoadPlugins []string
	SqlUri string
	Heroku bool
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

	returnConfig.ApiKey, ok = os.LookupEnv("BOT_API_KEY")// If env var is empty
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