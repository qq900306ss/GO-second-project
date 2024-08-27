package service

import (
	"fmt"
	"github/qq900306ss/SecondProject/moudle"
	"github/qq900306ss/SecondProject/utils"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// GetUserList
// @Summary 所有用戶
// @Tags 用戶資料
// @Success 200 {string} json{"code" , "message"}
// @Router /user/GetUserList [get]
func GetUserList(c *gin.Context) { //方法
	data := make([]*moudle.UserBasic, 10)
	data = moudle.GetUserList()

	c.JSON(http.StatusOK, gin.H{ //創造一個map json格式 gin.H就是回應
		"code":    0, // 0是成功 -1是失敗
		"message": "查詢成功",
		"data":    data,
	})

}

// CreateUser
// @Summary 新增用戶
// @Tags 用戶資料
// @Param name query string false "用戶名"
// @Param password query string false "密碼"
// @Param repassword query string false "確認密碼"
// @Success 200 {string} json{"code" , "message"}
// @Router /user/CreateUser [get]
func CreateUser(c *gin.Context) {
	user := moudle.UserBasic{}
	// user.Name = c.Query("name")
	// password := c.Query("password")
	// repassword := c.Query("repassword")

	user.Name = c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	repassword := c.Request.FormValue("repassword")
	fmt.Println("密碼: ", password)
	sweet := fmt.Sprintf("%06d", rand.Int31())

	data := moudle.FindUserByName(user.Name)
	if user.Name == "" || password == "" || repassword == "" {
		c.JSON(http.StatusOK, gin.H{ //創造一個map json格式 gin.H就是回應
			"code":    -1, // 0是成功 -1是失敗
			"message": "用戶名和帳號不可為空",
			"data":    user,
		})
		return

	}

	if data.Name != "" {
		c.JSON(http.StatusOK, gin.H{ //創造一個map json格式 gin.H就是回應
			"code":    -1, // 0是成功 -1是失敗
			"message": "用戶名已註冊",
			"data":    user,
		})
		return

	}

	if password != repassword {
		c.JSON(http.StatusOK, gin.H{ //創造一個map json格式 gin.H就是回應
			"code":    -1, // 0是成功 -1是失敗
			"message": "密碼不一樣",
			"data":    user,
		})
		return
	}

	// user.Password = password

	user.Password = utils.MakePassword(password, sweet) //加密密碼
	user.Sweet = sweet
	fmt.Println(user.Password)

	moudle.CreateUser(user)

	c.JSON(http.StatusOK, gin.H{ //創造一個map json格式 gin.H就是回應
		"code":    0, // 0是成功 -1是失敗
		"message": "創建成功",
		"data":    user,
	})

}

// FindUserByNameAndPwd
// @Summary 登入
// @Tags 用戶資料
// @Param name query string false "用戶名"
// @Param password query string false "密碼"
// @Success 200 {string} json{"code" , "message"}
// @Router /user/FindUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) { //方法
	data := moudle.UserBasic{}

	// name := c.Query("name")
	// password := c.Query("password")
	name := c.Request.FormValue("name")
	password := c.Request.FormValue("password")

	fmt.Println("有棟ㄇ?:", name, "password:", password)

	user := moudle.FindUserByName(name)
	fmt.Println("有東西?", user)

	if user.Name == "" {
		c.JSON(http.StatusOK, gin.H{ //創造一個map json格式 gin.H就是回應
			"code":    -1, // 0是成功 -1是失敗
			"message": "該用戶不存在",
			"data":    data,
		})
		return

	}

	flag := utils.ValidPassword(password, user.Sweet, user.Password) //驗證密碼正確性
	if !flag {
		c.JSON(http.StatusOK, gin.H{ //創造一個map json格式 gin.H就是回應
			"code":    -1, // 0是成功 -1是失敗
			"message": "密碼不正確捏",
			"data":    data,
		})
		return
	}
	pwd := utils.MakePassword(password, user.Sweet) //加密密碼

	data = moudle.FindUserByNameAndPwd(name, pwd) // 透過加密的密碼跟鳴子去找
	c.JSON(http.StatusOK, gin.H{                  //創造一個map json格式 gin.H就是回應
		"code":    0, // 0是成功 -1是失敗
		"message": "登入成功",
		"data":    data,
	})

}

// DeleteUser
// @Summary 刪除用戶
// @Tags 用戶資料
// @Param id query string false "id"
// @Success 200 {string} json{"code" , "message"}
// @Router /user/DeleteUser [get]
func DeleteUser(c *gin.Context) {
	user := moudle.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id")) //因為id 是個unint 所以要轉換成string
	user.ID = uint(id)                   //轉換成uint
	moudle.DeleteUser(user)

	c.JSON(http.StatusOK, gin.H{ //創造一個map json格式 gin.H就是回應
		"code":    0, // 0是成功 -1是失敗
		"message": "刪除成功",
		"data":    user,
	})

}

// UpdateUser
// @Summary 更新用戶
// @Tags 用戶資料
// @Param id formData string false "id"
// @Param name formData string false "name"
// @Param password formData string false "password"
// @Param phone formData string false "phone"
// @Param email formData string false "email"
// @Success 200 {string} json{"code" , "message"}
// @Router /user/UpdateUser [post]
func UpdateUser(c *gin.Context) {
	var user moudle.UserBasic
	// fmt.Println("測試測試:", c.PostForm("id"), "真假拉")   測試bug用
	// // c.Request.ParseForm()
	// // fmt.Println("Form data:", c.Request.Form) // 打印出所有表單數據

	idStr := c.PostForm("id")
	fmt.Println("Raw ID:", idStr) // 打印出 ID 的原始值
	id, err := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	if err != nil || id <= 0 { // 確保 ID 是有效的正整數
		c.JSON(http.StatusOK, gin.H{ //創造一個map json格式 gin.H就是回應
			"code":    -1, // 0是成功 -1是失敗
			"message": "ID錯誤",
			"data":    user,
		})
		return
	}

	// 更新用戶的字段
	user.Name = c.PostForm("name")
	user.Password = c.PostForm("password")
	user.Avatar = c.PostForm("icon") //頭像
	user.Email = c.PostForm("email")

	_, err1 := govalidator.ValidateStruct(user) // 驗證用戶資料

	if err1 != nil {
		fmt.Println(err1)
		c.JSON(http.StatusOK, gin.H{ //創造一個map json格式 gin.H就是回應
			"code":    -1, // 0是成功 -1是失敗
			"message": "更新失敗",
			"data":    user,
		})

	} else {
		c.JSON(http.StatusOK, gin.H{ //創造一個map json格式 gin.H就是回應
			"code":    0, // 0是成功 -1是失敗
			"message": "更新成功",
			"data":    user,
		})
		moudle.UpdateUser(user) // 更新用戶資料

	}

}

var upGrader = websocket.Upgrader{ //防止跨域站點偽造請求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(c, ws)
}

func RedisMsg(c *gin.Context) {
	userIdA, _ := strconv.Atoi(c.PostForm("userIdA"))
	userIdB, _ := strconv.Atoi(c.PostForm("userIdB"))
	start, _ := strconv.Atoi(c.PostForm("start")) //可以控制看幾次
	end, _ := strconv.Atoi(c.PostForm("end"))
	isRev, _ := strconv.ParseBool(c.PostForm("isRev")) //是否反轉
	res := moudle.RedisMsg(int64(userIdA), int64(userIdB), int64(start), int64(end), isRev)
	utils.RespOKList(c.Writer, "ok", res)
}

func MsgHandler(c *gin.Context, ws *websocket.Conn) {
	for {
		msg, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			fmt.Println(" MsgHandler 发送失败", err)
		}

		tm := time.Now().Format("2006-01-02 15:04:05")
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func SendUserMsg(c *gin.Context) { //目的：處理來自 HTTP 請求的 WebSocket 連接升級。
	moudle.Chat(c.Writer, c.Request)

}

func SearchFriend(c *gin.Context) {
	id, _ := strconv.Atoi(c.Request.FormValue("userId")) //因為id 是個unint 所以要轉換成string

	users := moudle.SearchFriend(uint(id))

	// c.JSON(http.StatusOK, gin.H{ //創造一個map json格式 gin.H就是回應 這裡返回這個不太好
	// 	"code":    0, // 0是成功 -1是失敗
	// 	"message": "查詢好友列表成功",
	// 	"data":    users,
	// })

	utils.RespOKList(c.Writer, users, len(users))

}

func AddFriend(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId")) //因為id 是個unint 所以要轉換成string
	targetId := c.Request.FormValue("targetId")

	code, msg := moudle.AddFriend(uint(userId), targetId)

	// c.JSON(http.StatusOK, gin.H{ //創造一個map json格式 gin.H就是回應 這裡返回這個不太好
	// 	"code":    0, // 0是成功 -1是失敗
	// 	"message": "查詢好友列表成功",
	// 	"data":    users,
	// })
	if code == 0 {
		utils.RespOK(c.Writer, code, msg)

	} else {
		utils.RespFail(c.Writer, msg)

	}

}

func CreateCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId")) //因為id 是個unint 所以要轉換成string
	Name := c.Request.FormValue("name")                        //因為id 是個unint 所以要轉換成string
	icon := c.Request.FormValue("icon")                        //因為id 是個unint 所以要轉換成string
	desc := c.Request.FormValue("desc")

	community := moudle.Community{}
	community.OwnerId = uint(ownerId)
	community.Name = Name
	community.Img = icon
	community.Desc = desc
	code, msg := moudle.CreateCommunity(community)

	moudle.JoinGroup(uint(ownerId), Name)

	if code == 0 {
		utils.RespOK(c.Writer, code, msg)

	} else {
		utils.RespFail(c.Writer, msg)

	}

}

func Loadcommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId")) //因為id 是個unint 所以要轉換成string

	data, msg := moudle.Loadcommunity(uint(ownerId))
	if len(data) != 0 {
		utils.RespList(c.Writer, 0, data, msg)

	} else {
		utils.RespFail(c.Writer, msg)

	}

}

func JoinGroup(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId")) //因為id 是個unint 所以要轉換成string
	comId := c.Request.FormValue("comId")                    //因為id 是個unint 所以要轉換成string

	data, msg := moudle.JoinGroup(uint(userId), comId)

	if data == 0 {
		utils.RespOK(c.Writer, data, msg)

	} else {
		utils.RespFail(c.Writer, msg)

	}

}
func FindByID(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))

	//	name := c.Request.FormValue("name")
	data := moudle.FindByID(uint(userId))
	utils.RespOK(c.Writer, data, "ok")
}
