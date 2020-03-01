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

package helpers

import (
	"fmt"
	"html"
	"strings"

	"github.com/HarukaNetwork/HarukaX/harukax/modules/sql"
	tg_md2html "github.com/PaulSonOfLars/gotg_md2html"
	"github.com/PaulSonOfLars/gotgbot/ext"
)

var MaxMessageLength = 4096

func MentionHtml(userId int, name string) string {
	return fmt.Sprintf("<a href=\"tg://user?id=%d\">%s</a>", userId, html.EscapeString(name))
}

func SplitMessage(msg string) []string {
	if len(msg) > MaxMessageLength {
		tmp := make([]string, 1)
		tmp[0] = msg
		return tmp
	} else {
		lines := strings.Split(msg, "\n")
		smallMsg := ""
		result := make([]string, 0)
		for _, line := range lines {
			if len(smallMsg)+len(line) < MaxMessageLength {
				smallMsg += line + "\n"
			} else {
				result = append(result, smallMsg)
				smallMsg = line + "\n"
			}
		}
		result = append(result, smallMsg)
		return result
	}
}

func SplitQuotes(text string) []string {
	smartOpen := "“"
	smartClose := "”"
	startChars := []string{"'", `"`, smartOpen}
	broke := false

	for _, char := range startChars {
		if strings.HasPrefix(text, char) {
			counter := 1
			for counter < len(text) {
				if text[counter] == '\\' {
					counter++
				} else if text[counter] == text[0] || (string(text[0]) == smartOpen && string(text[counter]) == smartClose) {
					broke = true
					break
				}
				counter++
			}
			if !broke {
				return strings.SplitN(text, " ", 2)
			}

			key := RemoveEscapes(strings.TrimSpace(text[1:counter]))
			rest := strings.TrimSpace(text[counter+1:])

			if key == "" {
				key = string(text[0]) + string(text[0])
			}
			tmp := make([]string, 2)
			tmp[0] = key
			tmp[1] = rest
			return tmp
		}
	}
	return strings.SplitN(text, " ", 2)
}

func RemoveEscapes(text string) string {
	counter := 0
	res := ""
	isEscaped := false

	for counter < len(text) {
		if isEscaped {
			res += string(text[counter])
			isEscaped = false
		} else if text[counter] == '\\' {
			isEscaped = true
		} else {
			res += string(text[counter])
		}
		counter++
	}
	return res
}

func BuildKeyboard(buttons []sql.Button) [][]ext.InlineKeyboardButton {
	keyb := make([][]ext.InlineKeyboardButton, 0)
	for _, btn := range buttons {
		if btn.SameLine && len(keyb) > 0 {
			keyb[len(keyb)-1] = append(keyb[len(keyb)-1], ext.InlineKeyboardButton{Text: btn.Name, Url: btn.Url})
		} else {
			k := make([]ext.InlineKeyboardButton, 1)
			k[0] = ext.InlineKeyboardButton{Text: btn.Name, Url: btn.Url}
			keyb = append(keyb, k)
		}
	}
	return keyb
}

func BuildWelcomeKeyboard(buttons []sql.WelcomeButton) [][]ext.InlineKeyboardButton {
	keyb := make([][]ext.InlineKeyboardButton, 0)
	for _, btn := range buttons {
		if btn.SameLine && len(keyb) > 0 {
			keyb[len(keyb)-1] = append(keyb[len(keyb)-1], ext.InlineKeyboardButton{Text: btn.Name, Url: btn.Url})
		} else {
			k := make([]ext.InlineKeyboardButton, 1)
			k[0] = ext.InlineKeyboardButton{Text: btn.Name, Url: btn.Url}
			keyb = append(keyb, k)
		}
	}
	return keyb
}

func GetNoteType(msg *ext.Message) (string, string, int, string, []tg_md2html.Button) {
	text := ""
	var dataType = -1
	var content string
	var rawText string
	var entities []ext.MessageEntity

	if reply := msg.ReplyToMessage; reply != nil {
		if reply.Text == "" {
			rawText = reply.Caption
			entities = reply.CaptionEntities
		} else {
			rawText = reply.Text
			entities = reply.Entities
		}
	} else {
		if msg.Text == "" {
			rawText = msg.Caption
			entities = msg.CaptionEntities
		} else {
			rawText = msg.Text
			entities = msg.Entities
		}
	}

	timesInserted := 0

	for _, ent := range entities {
		if ent.Type == "code" {
			rawText = rawText[:ent.Offset+timesInserted] + "`" + rawText[ent.Offset+timesInserted:]
			timesInserted++
			rawText = rawText[:(ent.Offset+ent.Length+(timesInserted))] + "`" + rawText[(ent.Offset+ent.Length+(timesInserted)):]
			timesInserted++
		}
	}

	args := strings.SplitN(msg.Text, " ", 3)
	noteName := args[1]

	buttons := make([]tg_md2html.Button, 0)
	if len(args) >= 3 {
		text, buttons = tg_md2html.MD2HTMLButtons(strings.SplitN(rawText, " ", 3)[2])

		if len(buttons) > 0 {
			dataType = sql.BUTTON_TEXT
		} else {
			dataType = sql.TEXT
		}
	} else if msg.ReplyToMessage != nil {
		if len(args) >= 2 && msg.ReplyToMessage.Text != "" {
			text, buttons = tg_md2html.MD2HTMLButtons(rawText)
			if len(buttons) > 0 {
				dataType = sql.BUTTON_TEXT
			} else {
				dataType = sql.TEXT
			}
		} else if msg.ReplyToMessage.Sticker != nil {
			content = msg.ReplyToMessage.Sticker.FileId
			dataType = sql.STICKER
		} else if msg.ReplyToMessage.Document != nil {
			content = msg.ReplyToMessage.Document.FileId
			text, buttons = tg_md2html.MD2HTMLButtons(rawText)
			dataType = sql.DOCUMENT
		} else if len(msg.ReplyToMessage.Photo) > 0 {
			content = msg.ReplyToMessage.Photo[len(msg.ReplyToMessage.Photo)-1].FileId
			text, buttons = tg_md2html.MD2HTMLButtons(rawText)
			dataType = sql.PHOTO
		} else if msg.ReplyToMessage.Audio != nil {
			content = msg.ReplyToMessage.Audio.FileId
			text, buttons = tg_md2html.MD2HTMLButtons(rawText)
			dataType = sql.AUDIO
		} else if msg.ReplyToMessage.Voice != nil {
			content = msg.ReplyToMessage.Voice.FileId
			text, buttons = tg_md2html.MD2HTMLButtons(rawText)
			dataType = sql.VOICE
		} else if msg.ReplyToMessage.Video != nil {
			content = msg.ReplyToMessage.Video.FileId
			text, buttons = tg_md2html.MD2HTMLButtons(rawText)
			dataType = sql.VIDEO
		}
	}
	return noteName, text, dataType, content, buttons
}

func GetWelcomeType(msg *ext.Message) (string, int, string, []tg_md2html.Button) {
	text := ""
	var dataType = -1
	var content string
	var rawText string
	var entities []ext.MessageEntity

	if reply := msg.ReplyToMessage; reply != nil {
		if reply.Text == "" {
			rawText = reply.Caption
			entities = reply.CaptionEntities
		} else {
			rawText = reply.Text
			entities = reply.Entities
		}
	} else {
		if msg.Text == "" {
			rawText = msg.Caption
			entities = msg.CaptionEntities
		} else {
			rawText = msg.Text
			entities = msg.Entities
		}
	}

	timesInserted := 0

	for _, ent := range entities {
		if ent.Type == "code" {
			rawText = rawText[:ent.Offset+timesInserted] + "`" + rawText[ent.Offset+timesInserted:]
			timesInserted++
			rawText = rawText[:(ent.Offset+ent.Length+(timesInserted))] + "`" + rawText[(ent.Offset+ent.Length+(timesInserted)):]
			timesInserted++
		}
	}

	args := strings.SplitN(msg.Text, " ", 2)

	buttons := make([]tg_md2html.Button, 0)
	if len(args) >= 2 {
		text, buttons = tg_md2html.MD2HTMLButtons(strings.SplitN(rawText, " ", 2)[1])

		if len(buttons) > 0 {
			dataType = sql.BUTTON_TEXT
		} else {
			dataType = sql.TEXT
		}
	} else if msg.ReplyToMessage != nil {
		if len(args) >= 2 && msg.ReplyToMessage.Text != "" {
			text, buttons = tg_md2html.MD2HTMLButtons(rawText)
			if len(buttons) > 0 {
				dataType = sql.BUTTON_TEXT
			} else {
				dataType = sql.TEXT
			}
		} else if msg.ReplyToMessage.Sticker != nil {
			content = msg.ReplyToMessage.Sticker.FileId
			dataType = sql.STICKER
		} else if msg.ReplyToMessage.Document != nil {
			content = msg.ReplyToMessage.Document.FileId
			text, buttons = tg_md2html.MD2HTMLButtons(rawText)
			dataType = sql.DOCUMENT
		} else if len(msg.ReplyToMessage.Photo) > 0 {
			content = msg.ReplyToMessage.Photo[len(msg.ReplyToMessage.Photo)-1].FileId
			text, buttons = tg_md2html.MD2HTMLButtons(rawText)
			dataType = sql.PHOTO
		} else if msg.ReplyToMessage.Audio != nil {
			content = msg.ReplyToMessage.Audio.FileId
			text, buttons = tg_md2html.MD2HTMLButtons(rawText)
			dataType = sql.AUDIO
		} else if msg.ReplyToMessage.Voice != nil {
			content = msg.ReplyToMessage.Voice.FileId
			text, buttons = tg_md2html.MD2HTMLButtons(rawText)
			dataType = sql.VOICE
		} else if msg.ReplyToMessage.Video != nil {
			content = msg.ReplyToMessage.Video.FileId
			text, buttons = tg_md2html.MD2HTMLButtons(rawText)
			dataType = sql.VIDEO
		}
	}
	return text, dataType, content, buttons
}

func RevertButtons(buttons []sql.WelcomeButton) string {
	res := ""
	for _, btn := range buttons {
		if btn.SameLine {
			res += fmt.Sprintf("\n[%s](buttonurl://%s)", btn.Name, btn.Url)
		} else {
			res += fmt.Sprintf("\n[%s](buttonurl://%s:same)", btn.Name, btn.Url)
		}
	}
	return res
}
