package service

import (
	"github/qq900306ss/SecondProject/moudle"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ResponseMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Get(c *gin.Context) {

	ind, err := template.ParseFiles("index.html", "views/chat/head.html")

	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "index")

}

// GetIndex
// @Tags 首頁
// @Success 200 {object} ResponseMessage // 使用自定義的響應結構
// @Router /index [get]
func GetIndex(c *gin.Context) {

	// ind, err := template.ParseFiles("index.html")
	// if err != nil {
	// 	panic(err)
	// }
	// ind.Execute(c.Writer, "index")

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Welcome !!",
	})
}

func ToRegister(c *gin.Context) {

	ind, err := template.ParseFiles("views/user/register.html")

	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "register")
}

func ToChant(c *gin.Context) {
	ind, err := template.ParseFiles("views/chat/index.html",
		"views/chat/head.html",
		"views/chat/tabmenu.html",
		"views/chat/concat.html",
		"views/chat/group.html",
		"views/chat/profile.html",
		"views/chat/createcom.html",
		"views/chat/userinfo.html",
		"views/chat/main.html",
		"views/chat/foot.html",
	)

	if err != nil {
		panic(err)
	}
	userId, _ := strconv.Atoi(c.Query("userId"))
	token := c.Query("token")
	user := moudle.UserBasic{}
	user.ID = uint(userId)
	user.Indentity = token
	ind.Execute(c.Writer, user)
}

func Chat(c *gin.Context) {
	moudle.Chat(c.Writer, c.Request)
}
