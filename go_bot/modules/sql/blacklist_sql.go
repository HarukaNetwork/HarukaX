package sql

type BlackListFilters struct {
	ChatId  string `gorm:"primary_key"`
	Trigger string `gorm:"primary_key"`
}

func AddToBlacklist(chatId string, trigger string) {
	filter := &BlackListFilters{ChatId: chatId, Trigger: trigger}
	SESSION.Save(filter)
}

func RmFromBlacklist(chatId string, trigger string) bool {
	filter := &BlackListFilters{ChatId: chatId, Trigger: trigger}
	if SESSION.Delete(filter).RowsAffected == 0 {
		return false
	}
	return true
}

func GetChatBlacklist(chatId string) []BlackListFilters {
	var filters []BlackListFilters
	SESSION.Where("chat_id = ?", chatId).Find(&filters)
	return filters
}