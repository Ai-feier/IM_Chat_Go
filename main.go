package main

import (
	"fmt"
	"ginchat/models"
	"ginchat/router"
	"ginchat/utils"
	"github.com/spf13/viper"
	"time"
)

func main() {
	utils.InitConfig()
	utils.InitMYSQL()
	utils.InitRedis()
	InitTimer()
	fmt.Println("MAIN:", utils.DB)
	r := router.Router()
	r.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// 初始化定时器
func InitTimer() {
	// 可优化：不是对每一个单独的节点进行清理，存在代码冗余，多次遍历所有在线节点，为减小系统负担可放长清理连接的时长
	utils.Timer(time.Duration(viper.GetInt("timeout.DelayHeartbeat"))*time.Second,
		time.Duration(viper.GetInt("timeout.HeartbeatHz"))*time.Second,
		models.CleanConnection, "")
}
