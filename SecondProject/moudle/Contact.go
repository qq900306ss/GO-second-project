package moudle

import (
	"fmt"
	"github/qq900306ss/SecondProject/utils"

	"gorm.io/gorm"
)

// 關係 關聯表
type Contact struct {
	gorm.Model
	OwnerId uint   //誰有的關係
	TagetId uint   //對應誰的關係
	Type    int    //關係類型 1 好友 2群組 3 廣播
	Desc    string //描述

}

func (table *Contact) TableName() string {
	return "contact"
}

func SearchFriend(userid uint) []UserBasic {
	contacts := make([]Contact, 0)
	objIds := make([]uint64, 0)
	utils.DB.Where("owner_id =? and type = 1", userid).Find(&contacts)
	for _, v := range contacts {
		fmt.Println(">>>>>>>>>", v)

		objIds = append(objIds, uint64(v.TagetId))
	}
	users := make([]UserBasic, 0)
	utils.DB.Where("id in ?", objIds).Find(&users)
	return users

}

func AddFriend(userId uint, targetName string) (int, string) {
	// user := UserBasic{}
	if targetName != "" {
		targetUser := FindUserByName(targetName)

		// fmt.Println(targetId, "      ", userId)
		if targetUser.Sweet != "" {
			if targetUser.ID == userId {
				return -1, "不能加自己"
			}
			contact0 := Contact{}

			utils.DB.Where("owner_id =? and taget_id =? and type = 1", userId, targetUser.ID).Find(&contact0)
			if contact0.ID != 0 {
				return -1, "已經是好友"
			}
			tx := utils.DB.Begin() //這裡用事務主要就是為了原子性(atomicity) 一個失敗那就全部回滾
			//事務一旦開始，不論期間甚麼異常最終都會rollback
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
				}
			}()

			contact := Contact{}
			contact.OwnerId = userId
			contact.TagetId = targetUser.ID
			contact.Type = 1
			contact.Desc = "好友"
			if err := utils.DB.Create(&contact).Error; err != nil {
				tx.Rollback()
				return -1, "加好友失敗"
			}
			contact1 := Contact{}
			contact1.OwnerId = targetUser.ID
			contact1.TagetId = userId
			contact1.Type = 1
			contact1.Desc = "好友"
			if err := utils.DB.Create(&contact1).Error; err != nil {
				tx.Rollback()
				return -1, "加好友失敗"
			}
			tx.Commit()
			return 0, "加好友成功"
		}
		return -1, "對方不存在"
	}
	return -1, "好友ID不能為空"
}
