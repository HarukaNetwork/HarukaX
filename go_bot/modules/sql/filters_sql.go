package sql

type CustomFilters struct {
	ChatId     string `sql:",pk"`
	Keyword    string `sql:",pk"`
	Reply      string
	IsSticker  bool `sql:",default:false"`
	IsDocument bool `sql:",default:false"`
	IsImage    bool `sql:",default:false"`
	IsAudio    bool `sql:",default:false"`
	IsVoice    bool `sql:",default:false"`
	IsVideo    bool `sql:",default:false"`
	HasButtons bool `sql:",default:false"`
}

type Buttons struct {
	Id int `sql:",pk"`
	ChatId string `sql:",pk"`
	Keyword string `sql:",pk"`
	Name string
	Url string
	SameLine bool `sql:",default:false"`
}

func AddFilter(chatId string, keyword string, reply string, isSticker bool, isDocument bool, isImage bool, isAudio bool, isVoice bool, isVideo bool, buttons []Buttons) {
	if buttons == nil {
		buttons = make([]Buttons, 0)
	}

}
