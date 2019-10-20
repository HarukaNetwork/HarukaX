package sql

import (
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/error_handling"
)

type Federation struct {
	Id        string `gorm:"primary_key"`
	OwnerId   string
	FedName   string
	FedAdmins []FedAdmin `gorm:"foreignkey:Id"`
	//FedChats []FedChat `gorm:"foreignkey:Id"`
	FedBans []FedBan `gorm:"foreignkey:FedRef"`
}

type FedAdmin struct {
	FedRef string `gorm:"primary_key"`
	UserId string `gorm:"primary_key"`
}

type FedChat struct {
	ChatId string `gorm:"primary_key"`
	FedRef string
}

type FedBan struct {
	FedRef string `gorm:"primary_key"`
	UserId string `gorm:"primary_key"`
	Reason string
}

func GetFedInfo(fedId string) *Federation {
	fed := &Federation{Id: fedId}
	if SESSION.First(fed).RowsAffected == 0 {
		return nil
	}
	return fed
} // No dirty reads

func GetFedFromOwnerId(ownerId string) *Federation {
	fed := &Federation{OwnerId: ownerId}
	if SESSION.Where("owner_id = ?", ownerId).First(fed).RowsAffected == 0 {
		return nil
	}
	return fed
} // No dirty reads

func GetFedId(chatId string) string {
	chat := &FedChat{}
	if SESSION.Where("chat_id = ?", chatId).First(chat).RowsAffected == 0 {
		return ""
	}
	return chat.FedRef
} // No dirty reads

func NewFed(ownerId string, fedId string, fedName string) bool {
	fed := &Federation{OwnerId: ownerId, Id: fedId, FedName: fedName}

	if err := SESSION.Save(fed).Error; err != nil {
		error_handling.HandleErr(err)
		return false
	}
	return true
} // No dirty read

func DelFed(fedId string) {
	tx := SESSION.Begin()

	fed := &Federation{}
	tx.Where("id = ?", fedId).Delete(fed)

	chat := &FedChat{}
	tx.Model(chat).Where("id = ?", fedId).Delete(chat)

	admins := &FedAdmin{}
	tx.Model(&admins).Where("id = ?", fedId).Delete(admins)

	bans := &FedBan{}
	tx.Model(bans).Where("id = ?", fedId).Delete(bans)

	tx.Commit()
} // No dirty reads

func IsUserFedAdmin(fedId string, userId string) string {
	fed := GetFedInfo(fedId)

	if fed.OwnerId == userId {
		return fed.OwnerId
	}

	admin := &FedAdmin{FedRef: fedId, UserId: userId}

	if SESSION.First(admin).RowsAffected == 0 {
		return ""
	} else {
		return admin.UserId
	}
} // No dirty reads

func GetChatFed(chatId string) *Federation {
	chat := &FedChat{ChatId: chatId}
	SESSION.Where("chat_id = ?", chatId).First(chat)
	return GetFedInfo(chat.FedRef)
} // No dirty reads

func ChatJoinFed(fedId string, chatId string) bool {
	chat := &FedChat{FedRef: fedId, ChatId: chatId}
	return SESSION.Save(chat).Error == nil
} // No dirty reads

func UserPromoteFed(fedId string, userId string) {
	admin := &FedAdmin{FedRef: fedId, UserId: userId}
	SESSION.Save(admin)
} //no dirty read

func UserDemoteFed(fedId string, userId string) {
	admin := &FedAdmin{FedRef: fedId, UserId: userId}
	SESSION.Delete(admin)
} // no dirty read

func ChatLeaveFed(chatId string) bool {
	chat := &FedChat{}
	return SESSION.Where("chat_id = ?", chatId).Delete(chat).RowsAffected != 0

} // no dirty read

func AllFedChats(fedId string) []string {
	var chats []FedChat
	SESSION.Where("fed_ref = ?", fedId).Find(&chats)
	tmp := make([]string, 0)
	for _, chat := range chats {
		tmp = append(tmp, chat.ChatId)
	}
	return tmp
} // no dirty read

func FbanUser(fedId string, userId string, reason string) {
	ban := &FedBan{FedRef: fedId, UserId: userId, Reason: reason}
	SESSION.Save(ban)
} // no dirty read

func UnFbanUser(fedId string, userId string) {
	ban := &FedBan{FedRef: fedId, UserId: userId}
	SESSION.Delete(ban)
} // no dirty read

func GetFbanUser(fedId string, userId string) *FedBan {
	ban := &FedBan{FedRef: fedId, UserId: userId}
	if SESSION.First(ban).RowsAffected == 0 {
		return nil
	}
	return ban
} // no dirty read

func GetAllFbanUsers(fedId string) []FedBan {
	var bans []FedBan
	SESSION.Where("id = ?", fedId).Find(&bans)
	return bans
} // no dirty read

func GetUserFbans(userId string) []Federation {
	var feds []Federation
	SESSION.Table("federations").Select("federations.id, federations.fed_name").
		Joins("left join fed_bans on fed_bans.fed_ref = federations.id").
		Where("fed_bans.user_id = ?", userId).Find(&feds)

	return feds
}

func GetAllFbanUsersGlobal() []FedBan {
	var bans []FedBan
	SESSION.Find(&bans)
	return bans
}

func GetAllFedsAdminsGlobal() []FedAdmin {
	var feds []FedAdmin
	SESSION.Find(&feds)
	return feds
}

func IsUserFedOwner(userId string, fedId string) bool {
	fed := GetFedInfo(fedId)
	return fed.OwnerId == userId
}

func GetFedAdmins(fedId string) []FedAdmin {
	var admins []FedAdmin
	SESSION.Where("fed_ref = ?", fedId).Find(&admins)
	return admins
}
