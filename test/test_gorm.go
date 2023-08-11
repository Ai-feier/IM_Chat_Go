package main

import (
	"ginchat/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	db, err := gorm.Open(mysql.Open(""), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	db.AutoMigrate(&models.UserBasic{})

	// Create
	// Create
	user := &models.UserBasic{
		Name:          "aifeier",
		LoginTime:     time.Now(),
		HeartbeatTime: time.Now(),
		LoginOutTime:  time.Now(),
		// 其他属性
	}
	db.Create(user)

	// Read
	db.First(user, 1) // 根据整型主键查找
	//db.First(&testUser, "code = ?", "D42") // 查找 code 字段值为 D42 的记录

	// Update - 将 product 的 price 更新为 200
	db.Model(user).Update("PassWord", 123456)
	//// Update - 更新多个字段
	//db.Model(&testUser).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	//db.Model(&testUser).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})
	//
	//// Delete - 删除 product
	//db.Delete(&testUser, 1)
}
