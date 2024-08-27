package moudle

import "gorm.io/gorm"

// 群信息表
type GroupBasic struct {
	gorm.Model
	Name    string
	OwnerId uint
	Icon    string
	Desc    string
	Type    int // 群類型

}

func (table *GroupBasic) TableName() string {
	return "group_basic"
}
