package models

import (
	"ginchat/utils"
	"gorm.io/gorm"
)

// 人员关系
type Contact struct {
	gorm.Model
	OwnerId  uint //谁的关系信息
	TargetId uint //对应的谁 /群 ID
	Type     int  //对应的类型  1好友  2群  3xx
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}

// AddFriend 添加好友   自己的ID  ， 好友的ID
func AddFriend(userId uint, targetName string) (int, string) {
	if targetName == "" {
		return -1, "好友ID不能为空"
	}
	targetUser := FindUserByName(targetName)
	if targetUser.Name == "" {
		return -1, "没有找到此用户"
	}
	if targetUser.ID == userId {
		return -1, "不能加自己"
	}
	// 查找关系表中是否已经存在当前关系
	var contact Contact
	utils.DB.Where("owner_id =?  and target_id =? and type=1", userId, targetUser.ID).Find(&contact)
	if contact.ID != 0 {
		return -1, "不能重复添加"
	}
	// 开启添加好友事务
	tx := utils.DB.Begin()
	// 报错就回滚事务
	defer func() {
		if r := recover(); r != nil {
			tx.Callback()
		}
	}()
	contact = Contact{
		OwnerId:  userId,
		TargetId: targetUser.ID,
		Type:     1,
	}
	err := utils.DB.Create(&contact).Error
	if err != nil {
		tx.Rollback()
		return -1, "添加好友失败"
	}
	contact = Contact{
		OwnerId:  targetUser.ID,
		TargetId: userId,
		Type:     1,
	}
	err = utils.DB.Create(&contact).Error
	if err != nil {
		tx.Rollback()
		return -1, "添加好友失败"
	}
	tx.Commit()
	return 0, "添加好友成功"
}
