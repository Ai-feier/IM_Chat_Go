package main

import (
	"ginchat/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	//db.AutoMigrate(&models.UserBasic{}) // user_basic表
	//db.AutoMigrate(&models.Contact{})    // contact表
	//db.AutoMigrate(&models.Message{})    // message表
	//db.AutoMigrate(&models.GroupBasic{}) // group_basic表
	db.AutoMigrate(&models.Community{})

	// Create
	// Create
	//user := &models.UserBasic{
	//	Name:          "aifeier",
	//	LoginTime:     time.Now(),
	//	HeartbeatTime: time.Now(),
	//	LoginOutTime:  time.Now(),
	//	// 其他属性
	//}
	//db.Create(user)
	//
	//// Read
	//db.First(user, 1) // 根据整型主键查找
	////db.First(&testUser, "code = ?", "D42") // 查找 code 字段值为 D42 的记录
	//
	//// Update - 将 product 的 price 更新为 200
	//db.Model(user).Update("PassWord", 123456)
	//// Update - 更新多个字段
	//db.Model(&testUser).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	//db.Model(&testUser).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})
	//
	//// Delete - 删除 product
	//db.Delete(&testUser, 1)
}
