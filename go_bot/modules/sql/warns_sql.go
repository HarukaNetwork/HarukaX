package sql

import (
	"fmt"
	"github.com/atechnohazard/ginko/go_bot/modules/utils/error_handling"
	"github.com/go-pg/pg/orm"
)

type Warns struct {
	UserId   string `sql:",pk"`
	ChatId   string `sql:",pk"`
	NumWarns int    `sql:",default:0"`
	Reasons  []string
}

func (w Warns) String() string {
	return fmt.Sprintf("<%d warns for %s in %s for reasons %s>", w.NumWarns, w.UserId, w.ChatId, w.Reasons)
}

type WarnFilters struct {
	ChatId  string `sql:",pk"`
	Keyword string `sql:"pk"`
	Reply   string `sql:",notnull"`
}

func (wf WarnFilters) String() string {
	return fmt.Sprintf("<Permissions for %d>", wf.ChatId)
}

type WarnSettings struct {
	ChatId    string `sql:",pk"`
	WarnLimit int    `sql:",default:3"`
	SoftWarn  bool   `sql:",default:false"`
}

func init() {
	models := []interface{}{&Warns{}, &WarnFilters{}, &WarnSettings{}}
	for _, model := range models {
		_ = SESSION.CreateTable(model, &orm.CreateTableOptions{FKConstraints: true})
	}
}

func WarnUser(userId string, chatId string, reason string) (int, []string) {
	warnedUser := &Warns{UserId: userId, ChatId: chatId}
	err := SESSION.Select(warnedUser)
	error_handling.HandleErrorGracefully(err)
	warnedUser.NumWarns++
	if reason != "" {
		warnedUser.Reasons = append(warnedUser.Reasons, reason)
	}
	reasons := warnedUser.Reasons
	num := warnedUser.NumWarns
	err = SESSION.Update(warnedUser)
	error_handling.HandleErrorGracefully(err)

	return num, reasons
}

func RemoveWarn(userId string, chatId string) bool {
	removed := false
	warnedUser := &Warns{UserId: userId, ChatId: chatId}
	err := SESSION.Select(warnedUser)

	if err == nil && warnedUser.NumWarns > 0 {
		warnedUser.NumWarns -= 1
		err := SESSION.Update(warnedUser)
		error_handling.HandleErrorGracefully(err)
		removed = true
	}

	return removed
}

func ResetWarns(userId string, chatId string) {
	warnedUser := &Warns{UserId: userId, ChatId: chatId}
	err := SESSION.Select(warnedUser)

	error_handling.HandleErrorGracefully(err)

	warnedUser.NumWarns = 0
	warnedUser.Reasons = make([]string, 0)
	err = SESSION.Update(warnedUser)
	error_handling.HandleErrorGracefully(err)
}

func GetWarns(userId string, chatId string) (int, []string) {
	user := &Warns{UserId: userId, ChatId: chatId}
	err := SESSION.Select(user)
	if err != nil {
		return 0, nil
	}
	return user.NumWarns, user.Reasons
}

func AddWarnFilter(chatId string, keyword string, reply string) {
	warnFilter := &WarnFilters{ChatId: chatId, Keyword: keyword, Reply: reply}
	err := SESSION.Insert(warnFilter)
	error_handling.HandleErrorGracefully(err)
}

func RemoveWarnFilter(chatId string, keyword string) bool {
	warnFilter := &WarnFilters{ChatId: chatId, Keyword: keyword}
	err := SESSION.Select(warnFilter)
	if err == nil {
		err := SESSION.Delete(warnFilter)
		error_handling.HandleErrorGracefully(err)
		return true
	}
	return false
}

func GetChatWarnTriggers(chatId string) *WarnFilters {
	warnFilters := &WarnFilters{ChatId: chatId}
	err := SESSION.Select(warnFilters)
	if err != nil {
		error_handling.HandleErrorGracefully(err)
		return nil
	} else {
		return warnFilters
	}
}

func GetWarnFilter(chatId string, keyword string) *WarnFilters {
	warnFilter := &WarnFilters{ChatId: chatId, Keyword: keyword}
	err := SESSION.Select(warnFilter)

	if err != nil {
		return nil
	} else {
		return warnFilter
	}
}

func SetWarnLimit(chatId string, warnLimit int) {
	warnSetting := &WarnSettings{ChatId: chatId}
	err := SESSION.Select(warnSetting)

	if err != nil {
		warnSetting = &WarnSettings{ChatId: chatId, WarnLimit: warnLimit}
	}
	warnSetting.WarnLimit = warnLimit

	err = SESSION.Update(warnSetting)
	error_handling.HandleErrorGracefully(err)
}

func SetWarnStrength(chatId string, softWarn bool) {
	warnSetting := &WarnSettings{ChatId: chatId}
	err := SESSION.Select(warnSetting)

	if err != nil {
		warnSetting = &WarnSettings{ChatId: chatId, SoftWarn: softWarn}
	}
	warnSetting.SoftWarn = softWarn

	err = SESSION.Update(warnSetting)
	error_handling.HandleErrorGracefully(err)
}

func GetWarnSetting(chatId string) (int, bool) {
	warnSetting := &WarnSettings{ChatId: chatId}
	err := SESSION.Select(warnSetting)
	if err != nil {
		return 3, false
	} else {
		return warnSetting.WarnLimit, warnSetting.SoftWarn
	}
}
