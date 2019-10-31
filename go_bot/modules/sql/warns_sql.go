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
	"fmt"
	"github.com/lib/pq"
)

type Warns struct {
	UserId   string         `gorm:"primary_key"`
	ChatId   string         `gorm:"primary_key"`
	NumWarns int            `gorm:"default:0"`
	Reasons  pq.StringArray `gorm:"type:varchar(64)[]"`
}

func (w Warns) String() string {
	return fmt.Sprintf("<%v warns for %s in %s for reasons %v>", w.NumWarns, w.UserId, w.ChatId, w.Reasons)
}

type WarnFilters struct {
	ChatId  string `gorm:"primary_key"`
	Keyword string `gorm:"primary_key"`
	Reply   string `gorm:"not null"`
}

func (wf WarnFilters) String() string {
	return fmt.Sprintf("<Permissions for %v>", wf.ChatId)
}

type WarnSettings struct {
	ChatId    string `gorm:"primary_key"`
	WarnLimit int    `gorm:"default:3"`
	SoftWarn  bool   `gorm:"default:false"`
}

func WarnUser(userId string, chatId string, reason string) (int, []string) {
	warnedUser := &Warns{UserId: userId, ChatId: chatId}
	SESSION.FirstOrInit(warnedUser)

	// Increment warns
	warnedUser.NumWarns++

	// Add reason if it exists
	if reason != "" {
		if len(reason) >= 64 {
			reason = reason[:63]
		}
		warnedUser.Reasons = append(warnedUser.Reasons, reason)
	}

	// Upsert warn
	SESSION.Save(warnedUser)

	return warnedUser.NumWarns, warnedUser.Reasons
}

func RemoveWarn(userId string, chatId string) bool {
	removed := false
	warnedUser := &Warns{UserId: userId, ChatId: chatId}
	SESSION.FirstOrInit(warnedUser)

	// only remove if user has warns
	if warnedUser.NumWarns > 0 {
		warnedUser.NumWarns -= 1
		SESSION.Save(warnedUser)
		removed = true
	}

	return removed
}

func ResetWarns(userId string, chatId string) {
	warnedUser := &Warns{UserId: userId, ChatId: chatId}
	SESSION.FirstOrInit(warnedUser)

	// resetting all warn fields
	warnedUser.NumWarns = 0
	warnedUser.Reasons = make([]string, 0)
	SESSION.Save(warnedUser)
}

func GetWarns(userId string, chatId string) (int, []string) {
	user := &Warns{UserId: userId, ChatId: chatId}
	SESSION.FirstOrInit(user)
	return user.NumWarns, user.Reasons
}

func AddWarnFilter(chatId string, keyword string, reply string) {
	warnFilter := &WarnFilters{ChatId: chatId, Keyword: keyword, Reply: reply}
	SESSION.Save(warnFilter)
}

func RemoveWarnFilter(chatId string, keyword string) bool {
	warnFilter := &WarnFilters{ChatId: chatId, Keyword: keyword}
	// return false if 0 rows were deleted
	if SESSION.Delete(warnFilter).RowsAffected == 0 {
		return false
	}
	return false
}

func GetChatWarnTriggers(chatId string) []WarnFilters {
	var warnFilters []WarnFilters
	SESSION.Where("chat_id = ?", chatId).Find(&warnFilters)
	if len(warnFilters) == 0 {
		return nil
	}
	return warnFilters
}

func GetWarnFilter(chatId string, keyword string) *WarnFilters {
	warnFilter := &WarnFilters{ChatId: chatId, Keyword: keyword}
	if SESSION.First(warnFilter).RowsAffected == 0 {
		return nil
	}
	return warnFilter
}

func SetWarnLimit(chatId string, warnLimit int) {
	warnSetting := &WarnSettings{ChatId: chatId}
	// init record if it doesn't exist
	SESSION.FirstOrInit(warnSetting)
	warnSetting.WarnLimit = warnLimit
	// upsert record
	SESSION.Save(warnSetting)
}

func SetWarnStrength(chatId string, softWarn bool) {
	warnSetting := &WarnSettings{ChatId: chatId, SoftWarn: softWarn}
	// init record if it doesn't exist
	SESSION.FirstOrInit(warnSetting)
	warnSetting.SoftWarn = softWarn
	// upsert record
	SESSION.Save(warnSetting)
}

func GetWarnSetting(chatId string) (int, bool) {
	warnSetting := &WarnSettings{ChatId: chatId}
	SESSION.FirstOrCreate(warnSetting)
	return warnSetting.WarnLimit, warnSetting.SoftWarn
}
