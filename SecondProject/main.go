package main

import (
	"github/qq900306ss/SecondProject/moudle"
	Router "github/qq900306ss/SecondProject/router"
	"github/qq900306ss/SecondProject/utils"
	"log"

	"github.com/spf13/viper"
)

func main() {
	utils.InitConfig() //初始化配置文件
	utils.InitMySQL()  //初始化連結database
	utils.InitRedis()  //初始化連結redis

	// 在啟動路由之前進行遷移
	err := utils.DB.AutoMigrate(
		&moudle.UserBasic{},
		&moudle.Message{},
		&moudle.Contact{},
		&moudle.GroupBasic{},
		&moudle.Community{},
	)
	if err != nil {
		log.Fatal("Database migration failed:", err)
	}
	r := Router.Router()
	r.Run(viper.GetString("port.server")) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
