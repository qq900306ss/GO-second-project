package moudle

import (
	"fmt"
	"github/qq900306ss/SecondProject/utils"
	"time"

	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name          string
	Password      string
	Phone         string `valid:"matches(^0[1-9]{1}\\d{8}$)"` //前面站兩個0開頭 第二個2~9 然後後面8個 ^表示開頭 $表示結尾
	Email         string `valid:"email"`
	Avatar        string
	Indentity     string
	ClentIP       string
	ClentPort     string
	Sweet         string //加密
	LoginTime     time.Time
	HeartbeatTime time.Time
	LoginOutTime  time.Time
	IsLogout      bool
	DeviceInfo    string
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic { //調用拿取資料的函式
	var data []*UserBasic //直接定義dat = 切片userbasic
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}
func FindUserByNameAndPwd(name string, password string) UserBasic { //登入檢查用

	user := UserBasic{}
	utils.DB.Where("name = ? and password = ?", name, password).First(&user) // First 只有一個ｆｉｎｄ集合

	str := fmt.Sprintf("%d", time.Now().Unix())
	temp := utils.MD5Encode(str) //taken加密
	fmt.Println(temp)
	utils.DB.Model(&user).Where("id = ?", user.ID).Update("indentity", temp)
	return user
}

func FindUserByName(name string) UserBasic {

	user := UserBasic{}
	utils.DB.Where("name = ?", name).First(&user) // First 只有一個ｆｉｎｄ集合
	return user
}

func FindUserByPhone(phone string) *gorm.DB {

	user := UserBasic{}
	return utils.DB.Where("phone = ?", phone).First(&user) // First 只有一個ｆｉｎｄ集合

}

func FindUserByEmail(email string) *gorm.DB {

	user := UserBasic{}
	return utils.DB.Where("email = ?", email).First(&user) // First 只有一個ｆｉｎｄ集合

}

func CreateUser(user UserBasic) *gorm.DB { //收到 返還一個gorm.db的指針
	user.LoginTime = time.Now() // 或指定其他有效的日期時間
	user.HeartbeatTime = time.Now()
	user.LoginOutTime = time.Now()
	return utils.DB.Create(&user)

}

func DeleteUser(user UserBasic) *gorm.DB { //刪除
	return utils.DB.Delete(&user)
}

func UpdateUser(user UserBasic) *gorm.DB { //更新
	return utils.DB.Model(&user).Updates(UserBasic{Name: user.Name, Password: user.Password, Phone: user.Phone, Email: user.Email, Avatar: user.Avatar})
}

func FindByID(id uint) UserBasic { //查詢id
	user := UserBasic{}
	utils.DB.Where("id = ?", id).First(&user) // First 只有一個ｆｉｎｄ集合
	return user

}
