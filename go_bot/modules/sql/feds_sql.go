package sql

import (
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/error_handling"
)

type Federations struct {
	OwnerId string `sql:",pk"`
	FedName string
	FedId   string
}

type ChatF struct {
	ChatId string `sql:",pk"`
	FedId  string
}

type UserF struct {
	UserId string `sql:",pk"`
	FedId  []string
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
	error_handling.HandleErr(err)
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

	var users []UserF
	err = SESSION.Model(&users).Select()
	error_handling.HandleErr(err)
	for _, user := range users {
		for i, fed := range user.FedId {
			if fedId == fed {
				user.FedId = append(user.FedId[:i], user.FedId[i+1:]...)
				if len(user.FedId) == 0 {
					_, err = SESSION.Model(&user).WherePK().Delete()
					error_handling.HandleErr(err)
				} else {
					_, err = SESSION.Model(user).OnConflict("(user_id) DO UPDATE").Set("fed_id = EXCLUDED.fed_id").Insert()
					error_handling.HandleErr(err)
				}
			}
		}
	}

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

func SearchUserInFed(fedId string, userId string) string {
	user := &UserF{}
	err := SESSION.Model(user).Where("user_id = ?", userId).Select()
	if err != nil {
		return ""
	}

	for _, fed := range user.FedId {
		if fed == fedId {
			return user.UserId
		}
	}

	return ""
}

func ChatJoinFed(fedId string, chatId string) bool {
	chat := &ChatF{FedId: fedId, ChatId: chatId}
	_, err := SESSION.Model(chat).OnConflict("(chat_id) DO UPDATE").Set("fed_id = EXCLUDED.fed_id").Insert()
	return err == nil
}

func UserDemoteFed(fedId string, userId string) bool {
	user := &UserF{UserId: userId}
	_ = SESSION.Model(user).WherePK().Select()

	for i, fed := range user.FedId {
		if fedId == fed {
			user.FedId = append(user.FedId[:i], user.FedId[i+1:]...)
			if len(user.FedId) == 0 {
				err := SESSION.Delete(user)
				return err == nil
			}
		}
	}

	_, err := SESSION.Model(user).OnConflict("(user_id) DO UPDATE").Set("fed_id = EXCLUDED.fed_id").Insert()
	return err == nil
}

func UserPromoteFed(fedId string, userId string) *UserF {
	user := &UserF{UserId: userId}
	_ = SESSION.Model(user).WherePK().Select()
	user.FedId = append(user.FedId, fedId)
	_, err := SESSION.Model(user).OnConflict("(user_id) DO UPDATE").Set("fed_id = EXCLUDED.fed_id").Insert()
	if err != nil {
		return nil
	}
	return user
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

func AllFedUsers(fedId string) []string {
	var users []UserF
	err := SESSION.Model(&users).Where("fed_id = ?", fedId).Select()
	error_handling.HandleErr(err)
	tmp := make([]string, 0)
	for _, user := range users {
		tmp = append(tmp, user.UserId)
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

func FbanUser(fedId string, userId string, reason string) *BansF {
	ban := &BansF{FedId: fedId, UserId: userId, Reason: reason}
	_, err := SESSION.Model(ban).OnConflict("(fed_id,user_id) DO UPDATE").Set("reason = EXCLUDED.reason").Insert()
	if err != nil {
		return nil
	} else {
		return ban
	}
}

func UnFbanUser(fedId string, userId string) *BansF {
	ban := &BansF{FedId: fedId, UserId: userId}
	_, err := SESSION.Model(ban).WherePK().Delete()
	if err != nil {
		return nil
	} else {
		return ban
	}
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

func GetAllFedsUsersGlobal() []string {
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
	} else {
		return fed
	}
}
