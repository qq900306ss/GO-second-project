package main

import (
	"github/qq900306ss/SecondProject/moudle"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:900306@tcp(127.0.0.1:3306)/secondproject?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{}) // parse time 把他轉go 的time.time 格式 loc 是當地時間
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&moudle.UserBasic{}) // DB生成表用 user的部分
	// db.AutoMigrate(&moudle.Message{}) // DB 生成表用 message的部分
	// db.AutoMigrate(&moudle.Contact{})    // DB 生成表用 contact的部分
	// db.AutoMigrate(&moudle.GroupBasic{}) // DB 生成表用 group的部分
	// db.AutoMigrate(&moudle.Community{}) // DB 生成表用 group的部分

	// Create
	// user := &moudle.UserBasic{}
	// user.Name = "帥啊"
	// user.LoginTime = time.Now()
	// user.HeartbeatTime = time.Now()
	// user.LoginOutTime = time.Now()

	// db.Create(user)

	// Read

	// fmt.Println(db.First(user, 1)) // find product with integer primary key
	// // db.First(&product, "code = ?", "D42") // find product with code D42

	// // Update - update product's price to 200
	// db.Model(user).Update("Password", "1234")
	// Update - update multiple fields
	// db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	// db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - delete product
	// db.Delete(&product, 1)
}
