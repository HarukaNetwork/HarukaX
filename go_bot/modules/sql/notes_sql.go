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

	"github.com/wI2L/jettison"

	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/caching"
)

const (
	TEXT        = 0
	BUTTON_TEXT = 1
	STICKER     = 2
	DOCUMENT    = 3
	PHOTO       = 4
	AUDIO       = 5
	VOICE       = 6
	VIDEO       = 7
)

type Note struct {
	ChatId     string `gorm:"primary_key" json:"chat_id"`
	Name       string `gorm:"primary_key" json:"name"`
	Value      string `gorm:"not null" json:"value"`
	File       string `json:"file"`
	IsReply    bool   `gorm:"default:false" json:"is_reply"`
	HasButtons bool   `gorm:"default:false" json:"has_buttons"`
	Msgtype    int    `gorm:"default:1" json:"msgtype"`
}

type Button struct {
	Id       uint   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	ChatId   string `gorm:"primary_key" json:"chat_id"`
	NoteName string `gorm:"primary_key" json:"note_name"`
	Name     string `gorm:"not null" json:"name"`
	Url      string `gorm:"not null" json:"url"`
	SameLine bool   `gorm:"default:false" json:"same_line"`
}

func AddNoteToDb(chatId string, noteName string, noteData string, msgtype int, buttons []Button, file string) {
	defer func() {
		cacheNote(chatId, noteName)
	}()
	if buttons == nil {
		buttons = make([]Button, 0)
	}

	tx := SESSION.Begin()

	prevButtons := make([]Button, 0)
	tx.Where(&Button{ChatId: chatId, NoteName: noteName}).Find(&prevButtons)
	for _, btn := range prevButtons {
		tx.Delete(&btn)
	}

	hasButtons := len(buttons) > 0

	note := &Note{ChatId: chatId, Name: noteName, Value: noteData, Msgtype: msgtype, File: file, HasButtons: hasButtons}
	tx.Where(Note{ChatId: chatId, Name: noteName}).Save(note)

	for _, btn := range buttons {
		btn := &Button{ChatId: chatId, NoteName: noteName, Name: btn.Name, Url: btn.Url, SameLine: btn.SameLine}
		tx.Create(btn)
	}
	tx.Commit()
}

func GetNote(chatId string, noteName string) *Note {
	notes, err := caching.CACHE.Get(fmt.Sprintf("note_%v", chatId))
	if err != nil {
		cacheNote(chatId, noteName)
	}

	var n []Note
	_ = json.Unmarshal(notes, &n)

	for _, note := range n {
		if note.Name == noteName {
			return &note
		}
	}

	return nil
}

func RmNote(chatId string, noteName string) bool {
	tx := SESSION.Begin()
	defer func() {
		cacheNote(chatId, noteName)
	}()
	note := &Note{ChatId: chatId, Name: noteName}

	if tx.First(note).RowsAffected == 0 {
		tx.Rollback()
		return false
	}

	buttons := make([]Button, 0)
	tx.Where(&Button{ChatId: chatId, NoteName: noteName}).Find(&buttons)
	for _, btn := range buttons {
		tx.Delete(&btn)
	}

	tx.Delete(note)
	tx.Commit()
	return true
}

func GetAllChatNotes(chatId string) []Note {
	notes, err := caching.CACHE.Get(fmt.Sprintf("note_%v", chatId))
	if err != nil {
		cacheNote(chatId, "")
	}

	var n []Note
	_ = json.Unmarshal(notes, &n)
	return n
}

func GetButtons(chatId string, noteName string) []Button {
	buttons, err := caching.CACHE.Get(fmt.Sprintf("button_%v_%v", chatId, noteName))
	if err != nil {
		cacheNote(chatId, noteName)
	}
	var btns []Button
	_ = json.Unmarshal(buttons, &btns)

	return btns
}

func cacheNote(chatId string, noteName string) {
	var notes []Note
	var buttons []Button
	if noteName != "" {
		SESSION.Where(Button{ChatId: chatId, NoteName: noteName}).Find(&buttons)
	}
	SESSION.Where("chat_id = ?", chatId).Find(&notes)

	if notes != nil {
		if len(notes) != 0 {
			nJson, _ := jettison.Marshal(notes)
			_ = caching.CACHE.Set(fmt.Sprintf("note_%v", chatId), nJson)
		}
	}

	if buttons != nil {
		if len(buttons) != 0 {
			nButtons, _ := jettison.Marshal(buttons)
			_ = caching.CACHE.Set(fmt.Sprintf("button_%v_%v", chatId, noteName), nButtons)
		}
	}
}
