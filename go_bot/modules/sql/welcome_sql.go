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
	CleanWelcome  int  `gorm:"default:0"`
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
	UserId        string `gorm:"primary_key"`
	ChatId        string `gorm:"primary_key"`
	ButtonClicked bool   `gorm:"default:false"`
}

// GetWelcomePrefs Return the preferences for welcoming users
func GetWelcomePrefs(chatID string) *Welcome {
	welc := &Welcome{ChatId: chatID}

	if SESSION.First(welc).RowsAffected == 0 {
		return &Welcome{
			ChatId:        chatID,
			ShouldWelcome: true,
			ShouldMute:    false,
			CleanWelcome:  0,
			DelJoined:     false,
			CustomWelcome: DefaultWelcome,
			WelcomeType:   TEXT,
			MuteTime:      0,
		}
	}
	return welc
}

// GetWelcomeButtons Get the buttons for the welcome message
func GetWelcomeButtons(chatID string) []WelcomeButton {
	var buttons []WelcomeButton
	SESSION.Where("chat_id = ?", chatID).Find(&buttons)
	return buttons
}

// SetCleanWelcome Set whether to clean old welcome messages or not
func SetCleanWelcome(chatID string, cw int) {
	w := &Welcome{ChatId: chatID, CleanWelcome: cw}
	SESSION.Save(w)
}

// UserClickedButton Mark the user as a human
func UserClickedButton(userID, chatID string) {
	mu := &MutedUser{UserId: userID, ChatId: chatID, ButtonClicked: true}
	SESSION.Save(mu)
}

// HasUserClickedButton Has the user clicked button to unmute themselves
func HasUserClickedButton(userID, chatID string) bool {
	mu := &MutedUser{UserId: userID, ChatId: chatID}
	SESSION.FirstOrInit(mu)
	return mu.ButtonClicked
}

// IsUserHuman Is the user a human
func IsUserHuman(userID, chatID string) bool {
	mu := &MutedUser{UserId: userID, ChatId: chatID}
	return SESSION.First(mu).RowsAffected != 0
}

func SetWelcPref(chatID string, pref bool) {
	w := &Welcome{ShouldWelcome:pref}
	SESSION.Save(w)
}