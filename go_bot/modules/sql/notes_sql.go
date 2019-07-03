package sql

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

func AddNoteToDb(chatId string, noteName string, noteData string, msgtype int, buttons []string, file string) {
	if buttons == nil {
		buttons = make([]string, 0)
	}
	note := &Note{ChatId: chatId, Name: noteName}
	SESSION.FirstOrCreate(note)
}
