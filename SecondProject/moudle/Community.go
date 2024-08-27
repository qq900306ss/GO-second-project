package moudle

import (
	"fmt"
	"github/qq900306ss/SecondProject/utils"

	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	Name    string
	OwnerId uint
	Img     string

	Desc string
}

func CreateCommunity(community Community) (int, string) {
	if len(community.Name) == 0 {
		return -1, "群組名稱不能為空"
	}
	if community.OwnerId == 0 {
		return -1, "請先登入"
	}
	if err := utils.DB.Create(&community).Error; err != nil {
		return -1, "建立群組失敗"
	}

	return 0, "建立群組成功"

}

func Loadcommunity(ownerId uint) ([]*Community, string) {
	objIds := make([]uint64, 0)
	var contact []*Contact

	utils.DB.Where("owner_id = ? and type = 2", ownerId).Find(&contact) //?不能跟and連載一起慧卡bug
	for _, v := range contact {
		objIds = append(objIds, uint64(v.TagetId))
	}
	data := make([]*Community, 10)

	utils.DB.Where("id in ?", objIds).Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}

	return data, "取得群組成功"
}
