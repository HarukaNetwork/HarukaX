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
	ChatId     string `gorm:"primary_key"`
	Name       string `gorm:"primary_key"`
	Value      string `gorm:"not null"`
	File       string
	IsReply    bool `gorm:"default:false"`
	HasButtons bool `gorm:"default:false"`
	Msgtype    int  `gorm:"default:1"`
}

type Button struct {
	Id       uint   `gorm:"primary_key;AUTO_INCREMENT"`
	ChatId   string `gorm:"primary_key"`
	NoteName string `gorm:"primary_key"`
	Name     string `gorm:"not null"`
	Url      string `gorm:"not null"`
	SameLine bool   `gorm:"default:false"`
}

func AddNoteToDb(chatId string, noteName string, noteData string, msgtype int, buttons []Button, file string) {
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
	note := &Note{ChatId: chatId, Name: noteName}
	if SESSION.First(note).RowsAffected == 0 {
		return nil
	}
	return note
}

func RmNote(chatId string, noteName string) bool {
	tx := SESSION.Begin()
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
	notes := make([]Note, 0)
	SESSION.Where("chat_id = ?", chatId).Find(&notes)
	return notes
}

func GetButtons(chatId string, noteName string) []Button {
	buttons := make([]Button, 0)
	if SESSION.Where(Button{ChatId: chatId, NoteName: noteName}).Find(&buttons).RowsAffected == 0 {
		return nil
	}
	return buttons
}
