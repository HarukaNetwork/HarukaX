package sql

import (
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/error_handling"
)

type Federations struct {
	OwnerId   string `sql:",pk"`
	FedName   string
	FedId     string
	FedAdmins []string
}

type ChatF struct {
	ChatId string `sql:",pk"`
	FedId  string
}

type RulesF struct {
	FedId string `sql:",pk"`
	Rules string
}

type BansF struct {
	FedId  string `sql:",pk"`
	UserId string `sql:",pk"`
	Reason string
}

func GetFedInfo(fedId string) *Federations {
	fed := &Federations{}
	err := SESSION.Model(fed).Where("fed_id = ?", fedId).Select()
	if err != nil {
		if err.Error() != "pg: no rows in result set" {
			error_handling.HandleErr(err)
		}
		return nil
	}
	return fed
}

func GetFedFromUser(userId string) *Federations {
	fed := &Federations{OwnerId: userId}
	err := SESSION.Model(fed).WherePK().Select()
	if err != nil {
		return nil
	}
	return fed
}

func GetFedId(chatId string) string {
	chat := &ChatF{}
	err := SESSION.Model(chat).Where("chat_id = ?", chatId).Select()
	if err != nil {
		return ""
	}
	return chat.FedId
}

func NewFed(ownerId string, fedId string, fedName string) string {
	fed := &Federations{OwnerId: ownerId, FedId: fedId, FedName: fedName}
	_, err := SESSION.Model(fed).OnConflict("(owner_id) DO UPDATE").Set("fed_name = EXCLUDED.fed_name").Insert()
	if err != nil {
		return ""
	}
	tmp := &Federations{OwnerId: ownerId}
	err = SESSION.Model(tmp).WherePK().Select()
	return tmp.FedId
}

func DelFed(fedId string) {
	fed := &Federations{}
	_, err := SESSION.Model(fed).Where("fed_id = ?", fedId).Delete()
	error_handling.HandleErr(err)

	chat := &ChatF{}
	_, err = SESSION.Model(chat).Where("fed_id = ?", fedId).Delete()
	error_handling.HandleErr(err)

	rules := &RulesF{}
	_, err = SESSION.Model(rules).Where("fed_id = ?", fedId).Delete()
	error_handling.HandleErr(err)
}

func SearchFedByName(fedName string) string {
	feds := &Federations{}
	err := SESSION.Model(feds).Where("fed_name = ?", fedName).Select()
	if err != nil {
		return ""
	} else {
		return feds.FedId
	}
}

func IsUserFedAdmin(fedId string, userId string) string {
	fed := GetFedInfo(fedId)
	if userId == fed.OwnerId {
		return userId
	}

	if len(fed.FedAdmins) == 0 {
		return ""
	}

	for _, user := range fed.FedAdmins {
		if userId == user {
			return user
		}
	}
	return ""
}

func GetChatFed(chatId string) *Federations {
	chat := &ChatF{ChatId: chatId}
	err := SESSION.Model(chat).WherePK().Select()
	if err != nil {
		return nil
	}
	return GetFedInfo(chat.FedId)
}

func ChatJoinFed(fedId string, chatId string) bool {
	chat := &ChatF{FedId: fedId, ChatId: chatId}
	_, err := SESSION.Model(chat).OnConflict("(chat_id) DO UPDATE").Set("fed_id = EXCLUDED.fed_id").Insert()
	return err == nil
}

func UserDemoteFed(fedId string, userId string) {
	federation := GetFedInfo(fedId)

	for i, fed := range federation.FedAdmins {
		if userId == fed {
			federation.FedAdmins = append(federation.FedAdmins[:i], federation.FedAdmins[i+1:]...)
		}
	}

	_, err := SESSION.Model(federation).OnConflict("(owner_id) DO UPDATE").Set("fed_admins = EXCLUDED.fed_admins").Insert()
	error_handling.HandleErr(err)
}

func UserPromoteFed(fedId string, userId string) {
	fed := GetFedInfo(fedId)
	fed.FedAdmins = append(fed.FedAdmins, userId)
	_, err := SESSION.Model(fed).OnConflict("(owner_id) DO UPDATE").Set("fed_admins = EXCLUDED.fed_admins").Insert()
	error_handling.HandleErr(err)
}

func ChatLeaveFed(chatId string) bool {
	chat := &ChatF{}
	_, err := SESSION.Model(chat).Where("chat_id = ?", chatId).Delete()
	return err == nil
}

func AllFedChats(fedId string) []string {
	var chats []ChatF
	err := SESSION.Model(&chats).Where("fed_id = ?", fedId).Select()
	error_handling.HandleErr(err)
	tmp := make([]string, 0)
	for _, chat := range chats {
		tmp = append(tmp, chat.ChatId)
	}
	return tmp
}

func SetFrules(fedId string, rules string) *RulesF {
	rule := &RulesF{FedId: fedId, Rules: rules}
	_, err := SESSION.Model(rule).OnConflict("(fed_id) DO UPDATE").Set("rules = EXCLUDED.rules").Insert()
	if err != nil {
		return nil
	} else {
		return rule
	}
}

func GetFrules(fedId string) *RulesF {
	rules := &RulesF{FedId: fedId}
	err := SESSION.Model(rules).WherePK().Select()
	if err != nil {
		return nil
	} else {
		return rules
	}
}

func FbanUser(fedId string, userId string, reason string) {
	ban := &BansF{FedId: fedId, UserId: userId, Reason: reason}
	_, err := SESSION.Model(ban).OnConflict("(fed_id,user_id) DO UPDATE").Set("reason = EXCLUDED.reason").Insert()
	error_handling.HandleErr(err)
}

func UnFbanUser(fedId string, userId string) {
	ban := &BansF{FedId: fedId, UserId: userId}
	_, err := SESSION.Model(ban).WherePK().Delete()
	error_handling.HandleErr(err)
}

func GetFbanUser(fedId string, userId string) *BansF {
	ban := &BansF{FedId: fedId, UserId: userId}
	err := SESSION.Model(ban).WherePK().Select()
	if err != nil {
		return nil
	} else {
		return ban
	}
}

func GetAllFbanUsers(fedId string) []BansF {
	var bans []BansF
	err := SESSION.Model(&bans).Where("fed_id = ?", fedId).Select()
	if err != nil {
		error_handling.HandleErr(err)
		return make([]BansF, 0)
	}

	return bans
}

func GetAllFbanUsersGlobal() []string {
	var bans []BansF
	err := SESSION.Model(&bans).Select()
	if err != nil {
		error_handling.HandleErr(err)
		return make([]string, 0)
	}

	tmp := make([]string, len(bans))
	for i, ban := range bans {
		tmp[i] = ban.UserId
	}
	return tmp
}

func GetAllFedsAdminsGlobal() []string {
	var feds []Federations
	err := SESSION.Model(&feds).Select()
	if err != nil {
		error_handling.HandleErr(err)
		return make([]string, 0)
	}

	tmp := make([]string, len(feds))
	for i, fed := range feds {
		tmp[i] = fed.FedId
	}
	return tmp
}

func IsUserFedOwner(userId string, fedId string) bool {
	fed := SearchFedById(fedId)
	return fed.OwnerId == userId
}

func SearchFedById(fedId string) *Federations {
	fed := &Federations{FedId: fedId}
	err := SESSION.Model(fed).Where("fed_id = ?", fedId).Select()
	if err != nil {
		return nil
	}
	return fed
}
