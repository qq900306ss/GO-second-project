package main

import (
	Router "github/qq900306ss/SecondProject/router"
	"github/qq900306ss/SecondProject/utils"

	"github.com/spf13/viper"
)

func main() {
	utils.InitConfig() //初始化配置文件
	utils.InitMySQL()  //初始化連結database
	utils.InitRedis()  //初始化連結redis
	r := Router.Router()
	r.Run(viper.GetString("port.server")) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
