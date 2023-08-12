package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	DB  *gorm.DB
	Red *redis.Client
)

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
}

func InitMYSQL() {
	// 配置日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	db, _ := gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{
		Logger: newLogger,
	})
	if db == nil {
		panic("db init false")
	}
	fmt.Println("MYSQL Inited !")
	DB = db
}

func InitRedis() {
	r := redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})
	Red = r
	fmt.Println("Redis Inited!")
}

const (
	PublishKey = "websocket"
)

// Publish 发布消息到Redis
func Publish(c context.Context, ch string, msg string) error {
	var err error
	fmt.Println("Publish...", msg)
	err = Red.Publish(c, ch, msg).Err()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// Subscribe 订阅Redis消息
func Subscribe(c context.Context, ch string) (string, error) {
	sub := Red.Subscribe(c, ch)
	fmt.Println("Subscribe...订阅消息")
	msg, err := sub.ReceiveMessage(c)
	if err != nil {
		fmt.Println("订阅Redis消息失败")
		return "", err
	}
	fmt.Println("Subscribe 。。。。", msg.Payload)
	return msg.String(), err
}
