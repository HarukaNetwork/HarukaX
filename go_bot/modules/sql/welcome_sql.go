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

const DefaultWelcome = "Hey {first}, how are you?"

type Welcome struct {
	ChatId        string `gorm:"primary_key"`
	CustomWelcome string
	ShouldWelcome bool `gorm:"default:true"`
	ShouldMute    bool `gorm:"default:true"`
	DelJoined     bool `gorm:"default:false"`
	CleanWelcome  int `gorm:"default:0"`
	WelcomeType   int  `gorm:"default:0"`
	MuteTime      int  `gorm:"default:0"`
}

type WelcomeButton struct {
	Id       uint   `gorm:"primary_key;AUTO_INCREMENT"`
	ChatId   string `gorm:"primary_key"`
	Name     string `gorm:"not null"`
	Url      string `gorm:"not null"`
	SameLine bool   `gorm:"default:false"`
}

type MutedUser struct {
	UserId string `gorm:"primary_key"`
	ChatId string `gorm:"primary_key"`
}

func GetWelcomePrefs(chatId string) *Welcome {
	welc := &Welcome{ChatId: chatId}

	if SESSION.First(welc).RowsAffected == 0 {
		return &Welcome{
			ChatId:        chatId,
			ShouldWelcome: true,
			ShouldMute:    true,
			CleanWelcome:  0,
			DelJoined:     false,
			CustomWelcome: DefaultWelcome,
			WelcomeType:   TEXT,
			MuteTime:      0,
		}
	}
	return welc
}

func GetWelcomeButtons(chatId string) []WelcomeButton {
	var buttons []WelcomeButton
	SESSION.Where("chat_id = ?", chatId).Find(&buttons)
	return buttons
}

func SetCleanWelcome(chatId string, cw int) {
	w := &Welcome{ChatId: chatId, CleanWelcome: cw}
	SESSION.Save(w)
}

func MarkUserHuman(userId string, chatId string) {
	mu := &MutedUser{UserId: userId, ChatId: chatId}
	SESSION.Save(mu)
}

func IsUserHuman(userId string, chatId string) bool {
	mu := &MutedUser{UserId: userId, ChatId: chatId}
	return SESSION.First(mu).RowsAffected != 0
}

//func SetCleanWelcome(chatId string, messageId int) {
//	w := &Welcome{ChatId: chatId, CleanWelcome: messageId}
//	SESSION.Save(w)
//}
