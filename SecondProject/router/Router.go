package Router

import (
	"github/qq900306ss/SecondProject/docs"
	"github/qq900306ss/SecondProject/service"

	"github.com/gin-gonic/gin"
	swaggerfile "github.com/swaggo/files"

	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine { //設置路由
	r := gin.Default() //創建一個新的gin實例
	//swager設置

	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfile.Handler)) //web框架加上swagger 設備產生的api

	//靜態資料
	r.Static("/asset", "asset/") //設置靜態資料夾路徑
	r.LoadHTMLGlob("views/**/*") //設置模板路徑下面所有的

	//首頁
	r.GET("/", service.Get)                  //用戶/index之後會調用service.GetIndex這函數來處理
	r.GET("/index", service.GetIndex)        //用戶/index之後會調用service.GetIndex這函數來處理
	r.GET("/ToRegister", service.ToRegister) //用戶/index之後會調用service.GetIndex這函數來處理
	r.GET("/ToChant", service.ToChant)       //用戶/index之後會調用service.GetIndex這函數來處理
	r.GET("/Chat", service.Chat)             //用戶/index之後會調用service.GetIndex這函數來處理

	r.POST("/SearchFriend", service.SearchFriend)

	//用戶管理
	r.POST("/user/GetUserList", service.GetUserList)
	r.POST("/user/CreateUser", service.CreateUser)
	r.POST("/user/DeleteUser", service.DeleteUser)
	r.POST("/user/FindUserByNameAndPwd", service.FindUserByNameAndPwd)
	r.POST("/user/UpdateUser", service.UpdateUser)
	r.POST("/user/find", service.FindByID)

	//發送消息
	r.GET("/user/SendMsg", service.SendMsg)
	r.GET("/user/SendUserMsg", service.SendUserMsg)

	r.POST("/attach/upload", service.Upload)
	r.POST("/contact/AddFriend", service.AddFriend)

	r.POST("/contact/CreateCommunity", service.CreateCommunity)

	r.POST("/contact/Loadcommunity", service.Loadcommunity)

	r.POST("/contact/JoinGroup", service.JoinGroup)
	r.POST("/user/RedisMsg", service.RedisMsg)

	return r

}
