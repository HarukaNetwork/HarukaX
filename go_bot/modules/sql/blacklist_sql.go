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
	"encoding/json"
	"fmt"

	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/caching"
)

type BlackListFilters struct {
	ChatID  string `gorm:"primary_key" json:"chat_id"`
	Trigger string `gorm:"primary_key" json:"trigger"`
}

func AddToBlacklist(chatID string, trigger string) {
	filter := &BlackListFilters{ChatID: chatID, Trigger: trigger}
	SESSION.Save(filter)
	go CacheBlacklist(chatID)
}

func RmFromBlacklist(chatID string, trigger string) bool {
	filter := &BlackListFilters{ChatID: chatID, Trigger: trigger}
	defer func(chatID string) {
		go CacheBlacklist(chatID)
	}(chatID)
	if SESSION.Delete(filter).RowsAffected == 0 {
		return false
	}
	return true
}

func GetChatBlacklist(chatID string) []BlackListFilters {
	blf, err := caching.CACHE.Get(fmt.Sprintf("blacklist_%v", chatID))
	if err != nil {
		go CacheBlacklist(chatID)
	}

	var blistFilters []BlackListFilters
	_ = json.Unmarshal(blf, &blistFilters)
	return blistFilters
}

func CacheBlacklist(chatID string) {
	var filters []BlackListFilters
	SESSION.Where("chat_id = ?", chatID).Find(&filters)
	blJson, _ := json.Marshal(filters)
	_ = caching.CACHE.Set(fmt.Sprintf("blacklist_%v", chatID), blJson)
}
