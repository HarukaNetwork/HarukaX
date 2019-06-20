package sql

import (
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/error_handling"
)

type Federations struct {
	OwnerId string
	FedName string
	FedId   string `sql:",pk"`
}

type FedAdmins struct {
	FedId  string `sql:",pk" pg:"fk:fed_id"`
	UserId string `sql:",pk"`
}

type FedChats struct {
	ChatId string `sql:",pk"`
	FedId  string `pg:"fk:fed_id"`
}

type FedBans struct {
	FedId  string `sql:",pk" pg:"fk:fed_id"`
	UserId string `sql:",pk"`
	Reason string
}

func GetFedInfo(fedId string) *Federations {
	fed := &Federations{FedId: fedId}
	err := SESSION.Model(fed).Where("fed_id = ?", fedId).Select()
	if err != nil {
		if err.Error() != "pg: no rows in result set" {
			error_handling.HandleErr(err)
		}
		return nil
	}
	return fed
} // No dirty reads

func GetFedFromOwnerId(ownerId string) *Federations {
	fed := &Federations{OwnerId: ownerId}
	err := SESSION.Model(fed).Where("owner_id = ?", ownerId).Select()
	if err != nil {
		return nil
	}
	return fed
} // No dirty reads

func GetFedId(chatId string) string {
	chat := &FedChats{}
	err := SESSION.Model(chat).Where("chat_id = ?", chatId).Select()
	if err != nil {
		return ""
	}
	return chat.FedId
} // No dirty reads

func NewFed(ownerId string, fedId string, fedName string) bool {
	fed := &Federations{OwnerId: ownerId, FedId: fedId, FedName: fedName}
	_, err := SESSION.Model(fed).OnConflict("(fed_id) DO UPDATE").Set("fed_name = EXCLUDED.fed_name").Insert()
	if err != nil {
		error_handling.HandleErr(err)
		return false
	}
	return true
} // No dirty read

func DelFed(fedId string) {
	fed := &Federations{}
	_, err := SESSION.Model(fed).Where("fed_id = ?", fedId).Delete()
	error_handling.HandleErr(err)

	chat := &FedChats{}
	_, err = SESSION.Model(chat).Where("fed_id = ?", fedId).Delete()
	error_handling.HandleErr(err)

	var admins []FedAdmins
	_, err = SESSION.Model(&admins).Where("fed_id = ?", fedId).Delete()

	bans := &FedBans{}
	_, err = SESSION.Model(bans).Where("fed_id = ?", fedId).Delete()
	error_handling.HandleErr(err)
} // No dirty reads

func IsUserFedAdmin(fedId string, userId string) string {
	admin := &FedAdmins{FedId: fedId, UserId: userId}
	err := SESSION.Select(admin)
	if err != nil {
		if err.Error() != "pg: no rows in result set" {
			error_handling.HandleErr(err)
		}
		return ""
	} else {
		return admin.UserId
	}
} // No dirty reads

func GetChatFed(chatId string) *Federations {
	chat := &FedChats{ChatId: chatId}
	err := SESSION.Model(chat).Where("chat_id = ?", chatId).Select()
	if err != nil {
		return nil
	}
	return GetFedInfo(chat.FedId)
} // No dirty reads

func ChatJoinFed(fedId string, chatId string) bool {
	chat := &FedChats{FedId: fedId, ChatId: chatId}
	_, err := SESSION.Model(chat).OnConflict("(chat_id) DO UPDATE").Set("fed_id = EXCLUDED.fed_id").Insert()
	return err == nil
} // No dirty reads

func UserPromoteFed(fedId string, userId string) {
	admin := &FedAdmins{FedId: fedId, UserId: userId}
	err := SESSION.Insert(admin)
	error_handling.HandleErr(err)
} //no dirty read

func UserDemoteFed(fedId string, userId string) {
	admin := &FedAdmins{FedId: fedId, UserId: userId}
	err := SESSION.Delete(admin)
	error_handling.HandleErr(err)
} // no dirty read

func ChatLeaveFed(chatId string) bool {
	chat := &FedChats{}
	_, err := SESSION.Model(chat).Where("chat_id = ?", chatId).Delete()
	return err == nil
} // no dirty read

func AllFedChats(fedId string) []string {
	var chats []FedChats
	err := SESSION.Model(&chats).Where("fed_id = ?", fedId).Select()
	if err != nil {
		if err.Error() != "pg: no rows in result set" {
			error_handling.HandleErr(err)
		}
	}
	tmp := make([]string, 0)
	for _, chat := range chats {
		tmp = append(tmp, chat.ChatId)
	}
	return tmp
} // no dirty read

func FbanUser(fedId string, userId string, reason string) {
	ban := &FedBans{FedId: fedId, UserId: userId, Reason: reason}
	_, err := SESSION.Model(ban).OnConflict("(fed_id,user_id) DO UPDATE").Set("reason = EXCLUDED.reason").Insert()
	error_handling.HandleErr(err)
} // no dirty read

func UnFbanUser(fedId string, userId string) {
	ban := &FedBans{FedId: fedId, UserId: userId}
	_, err := SESSION.Model(ban).WherePK().Delete()
	error_handling.HandleErr(err)
} // no dirty read

func GetFbanUser(fedId string, userId string) *FedBans {
	ban := &FedBans{FedId: fedId, UserId: userId}
	err := SESSION.Model(ban).WherePK().Select()
	if err != nil {
		if err.Error() != "pg: no rows in result set" {
			error_handling.HandleErr(err)
		}
		return nil
	} else {
		return ban
	}
} // no dirty read

func GetAllFbanUsers(fedId string) []FedBans {
	var bans []FedBans
	err := SESSION.Model(&bans).Where("fed_id = ?", fedId).Select()
	if err != nil {
		if err.Error() != "pg: no rows in result set" {
			error_handling.HandleErr(err)
		}
		return make([]FedBans, 0)
	}
	return bans
} // no dirty read

func GetAllFbanUsersGlobal() []string {
	var bans []FedBans
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
	fed := GetFedInfo(fedId)
	return fed.OwnerId == userId
}

func GetFedAdmins(fedId string) []FedAdmins {
	var admins []FedAdmins
	err := SESSION.Model(&admins).Where("fed_id = ?", fedId).Select()
	if err != nil {
		if err.Error() != "pg: no rows in result set" {
			error_handling.HandleErr(err)
		}
		return make([]FedAdmins, 0)
	} else {
		return admins
	}
}
