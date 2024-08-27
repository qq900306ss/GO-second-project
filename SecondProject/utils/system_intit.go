package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitConfig() { //讀取yml的sql位置
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Config file loaded app:", viper.Get("app"))
	fmt.Println("Config file loaded mysql", viper.Get("mysql"))

}

var (
	DB  *gorm.DB
	Red *redis.Client
)

func InitMySQL() { //連接mysql
	//自訂日誌模版 打印sql
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢sql 的閥值
			LogLevel:      logger.Info, // 等級
			Colorful:      true,        // 彩色化輸出
		},
	)
	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{Logger: newLogger}) // parse time 把他轉go 的time.time 格式 loc 是當地時間
	// user := &moudle.UserBasic{}
	// DB.Find(&user)
	// fmt.Println(user)

}

func InitRedis() { //連接redis
	Red = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConns"),
	})

	pong, err := Red.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("redis連線有問題", err)
	} else {
		fmt.Println("redis連線成功", pong)
	}

}

const (
	PublishKey = "websocket"
)

func Publish(ctx context.Context, channel string, msg string) error { //發布消息到redis 目的：將消息發布到指定的 Redis 頻道。

	var err error
	fmt.Println("Publish : ", msg)

	err = Red.Publish(ctx, channel, msg).Err()
	if err != nil {
		fmt.Println("發布redis消息有問題", err)
		return err

	}
	return err
}

func Subscribe(ctx context.Context, channel string) (string, error) { //訂閱消息到Redis 目的：訂閱指定的 Redis 頻道，並接收消息。

	sub := Red.Subscribe(ctx, channel)
	fmt.Println("Subscribe ctx : ", ctx) // 終端顯示一下有沒有收到 //都是指針
	msg, err := sub.ReceiveMessage(ctx)  // 如果沒收到消息就會一直等待

	if err != nil {
		fmt.Println("訂閱redis消息有問題", err)
		return "", err

	}
	fmt.Println("Subscribe payload : ", msg.Payload)

	return msg.Payload, err

}
