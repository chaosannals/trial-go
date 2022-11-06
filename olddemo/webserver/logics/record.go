package logics

import (
	"fmt"
	"github.com/chaosannals/trial-go/models"
	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
	"gorm.io/gorm"
	"github.com/satori/go.uuid"
)

var db *gorm.DB = nil
var title *riot.Engine = nil
var content *riot.Engine = nil

//Init 初始化
func Init() {
	db = models.InitStore()
	title = models.AdpotIndex("title")
	content = models.AdpotIndex("content")
}

//Recover 回收资源
func Recover() {
	title.Close()
	content.Close()
}

//Insert 插入
func Insert(data models.Store) *gorm.DB {
	if id, err := uuid.NewV1(); err != nil {

	} else {
		data.ID = id.String()
	}
	title.Index(data.ID, types.DocData{Content: data.Title}, false)
	title.Flush()
	content.Index(data.ID, types.DocData{Content: data.Content}, false)
	content.Flush()
	return db.Create(&data)
}

//Update 修改
func Update(data models.Store) *gorm.DB {
	title.Index(data.ID, types.DocData{Content: data.Title}, true)
	title.Flush()
	content.Index(data.ID, types.DocData{Content: data.Content}, true)
	content.Flush()
	return db.Save(&data)
}

//Remove 删除
func Remove(id string) {
	db.Unscoped().Where(fmt.Sprintf("ID = %s", id)).Delete(&models.Store{})
	title.RemoveDoc(id, true)
	title.Flush()
	content.RemoveDoc(id, true)
	content.Flush()
}

//Search 搜索
func Search(request types.SearchReq) types.SearchResp {
	return content.Search(request)
}
