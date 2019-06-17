package sql

import (
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/error_handling"
)

type BlackListFilters struct {
	ChatId  string `sql:",pk"`
	Trigger string `sql:",pk"`
}


func AddToBlacklist(chatId string, trigger string) {
	filter := &BlackListFilters{ChatId: chatId, Trigger: trigger}
	_ = SESSION.Insert(filter)
}

func RmFromBlacklist(chatId string, trigger string) bool {
	filter := &BlackListFilters{ChatId: chatId, Trigger: trigger}
	err := SESSION.Select(filter)
	if err != nil {
		return false
	} else {
		err := SESSION.Delete(filter)
		error_handling.HandleErr(err)
		return true
	}
}

func GetChatBlacklist(chatId string) []BlackListFilters {
	var filters []BlackListFilters
	err := SESSION.Model(&filters).Where("chat_id = ?", chatId).Select()
	error_handling.HandleErr(err)
	return filters
}

func NumBlacklistFilters(chatId string) int {
	count, err := SESSION.Model((*BlackListFilters)(nil)).Where("chat_id = ?", chatId).Count()
	error_handling.HandleErr(err)
	return count
}