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
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string `valid:"email"`
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
	utils.DB.Table(&UserBasic{}).Joins()
	return data
}

// FindByNameAndPwd 通过用户名和密码查询
func FindByNameAndPwd(name, password string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ? and pass_word = ?", name, password).First(&user)
	return user
}

// 通过用户名查找
func FindByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ?", name).First(&user)
	return user
}

// 通过手机号查找
func FindByPhone(phone string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("phone = ?", phone).First(&user)
	return user
}

// 通过邮箱查找
func FindByEmail(email string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("email = ?", email).First(&user)
	return user
}

// CreateUser 创建用户
func CreateUser(user *UserBasic) *gorm.DB {
	return utils.DB.Create(user)
}

// DeleteUser 删除用户
func DeleteUser(user *UserBasic) *gorm.DB {
	return utils.DB.Delete(user)
}

// UpdateUser 修改用户
func UpdateUser(user *UserBasic) *gorm.DB {
	return utils.DB.Model(user).Updates(UserBasic{Name: user.Name, PassWord: user.PassWord, Phone: user.Phone, Email: user.Email, Avatar: user.Avatar})
}
