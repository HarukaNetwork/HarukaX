package error_handling

import (
	"log"
)

type CommandCallback func()

func HandleErr(err error) {
	if err != nil {
		log.Println("Error: ", err)
	}
}

func HandleErrorAndExit(err error) {
	if err != nil {
		log.Fatal("Error: ", err)
	}
}