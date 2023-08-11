package models

import (
	"fmt"
	"ginchat/utils"
	"time"

	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name          string
	PassWord      string
	Phone         string
	Email         string
	Avatar        string //头像
	Identity      string
	ClientIp      string
	ClientPort    string
	Salt          string
	LoginTime     time.Time
	HeartbeatTime time.Time
	LoginOutTime  time.Time `gorm:"column:login_out_time" json:"login_out_time"`
	IsLogout      bool
	DeviceInfo    string
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 8)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

// CreateUser 创建用户
func CreateUser(user *UserBasic) *gorm.DB {
	return utils.DB.Create(user)
}

// DeleteUser 删除用户
func DeleteUser(user *UserBasic) {
	utils.DB.Delete(user)
}

// UpdateUser 修改用户
func UpdateUser(user *UserBasic) *gorm.DB {
	return utils.DB.Model(user).Updates(UserBasic{Name: user.Name, PassWord: user.PassWord, Phone: user.Phone, Email: user.Email, Avatar: user.Avatar})
}
