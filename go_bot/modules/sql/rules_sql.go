package sql

import (
	"encoding/json"
	"fmt"

	"github.com/wI2L/jettison"

	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/caching"
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/error_handling"
)

type Rules struct {
	ChatId string `gorm:"primary_key" json:"chat_id"`
	Rules  string `json:"rules"`
}

func GetChatRules(chatId string) *Rules {
	ruleJson, err := caching.CACHE.Get(fmt.Sprintf("rules_%v", chatId))
	var rules *Rules
	if err != nil {
		rules = cacheRules(chatId)
	}
	_ = json.Unmarshal(ruleJson, &rules)
	return rules
}

func SetChatRules(chatId, rules string) {
	SESSION.Save(&Rules{ChatId: chatId, Rules: rules})
	cacheRules(chatId)
}

func cacheRules(chatId string) *Rules {
	rules := &Rules{}
	SESSION.Where("chat_id = ?", chatId).Find(&rules)
	ruleJson, _ := jettison.Marshal(&rules)
	err := caching.CACHE.Set(fmt.Sprintf("rules_%v", chatId), ruleJson)
	error_handling.HandleErr(err)
	return rules
}
