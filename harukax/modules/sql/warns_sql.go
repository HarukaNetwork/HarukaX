/*
 *    Copyright © 2020 Haruka Network Development
 *    This file is part of Haruka X.
 *
 *    Haruka X is free software: you can redistribute it and/or modify
 *    it under the terms of the Raphielscape Public License as published by
 *    the Devscapes Open Source Holding GmbH., version 1.d
 *
 *    Haruka X is distributed in the hope that it will be useful,
 *    but WITHOUT ANY WARRANTY; without even the implied warranty of
 *    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *    Devscapes Raphielscape Public License for more details.
 *
 *    You should have received a copy of the Devscapes Raphielscape Public License
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
	tx := SESSION.Begin()
	tx.FirstOrInit(warnedUser)

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
	tx.Save(warnedUser)
	tx.Commit()

	return warnedUser.NumWarns, warnedUser.Reasons
}

func RemoveWarn(userId string, chatId string) bool {
	removed := false
	warnedUser := &Warns{UserId: userId, ChatId: chatId}
	tx := SESSION.Begin()

	tx.FirstOrInit(warnedUser)

	// only remove if user has warns
	if warnedUser.NumWarns > 0 {
		warnedUser.NumWarns -= 1
		tx.Save(warnedUser)
		removed = true
	}
	tx.Commit()

	return removed
}

func ResetWarns(userId string, chatId string) {
	warnedUser := &Warns{UserId: userId, ChatId: chatId}
	tx := SESSION.Begin()

	tx.FirstOrInit(warnedUser)

	// resetting all warn fields
	warnedUser.NumWarns = 0
	warnedUser.Reasons = make([]string, 0)
	tx.Save(warnedUser)
	tx.Commit()
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
	tx := SESSION.Begin()
	// init record if it doesn't exist
	tx.FirstOrInit(warnSetting)
	warnSetting.WarnLimit = warnLimit
	// upsert record
	tx.Save(warnSetting)
}

func SetWarnStrength(chatId string, softWarn bool) {
	warnSetting := &WarnSettings{ChatId: chatId, SoftWarn: softWarn}
	tx := SESSION.Begin()
	// init record if it doesn't exist
	tx.FirstOrInit(warnSetting)
	warnSetting.SoftWarn = softWarn
	// upsert record
	tx.Save(warnSetting)
	tx.Commit()
}

func GetWarnSetting(chatId string) (int, bool) {
	warnSetting := &WarnSettings{ChatId: chatId}
	SESSION.FirstOrCreate(warnSetting)
	return warnSetting.WarnLimit, warnSetting.SoftWarn
}
