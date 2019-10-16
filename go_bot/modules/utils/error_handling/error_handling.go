package error_handling

import (
	log "github.com/sirupsen/logrus"
)

type CommandCallback func()

func HandleErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func FatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
