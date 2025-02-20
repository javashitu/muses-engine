package repo

import (
	"log"
	"muses-engine/config"
	"muses-engine/internal/model/entity"
	"muses-engine/pkg/util"

	"github.com/bwmarrin/snowflake"
)

func SaveFileRecord(fileStore entity.FileStore) (*entity.FileStore, error) {
	if fileStore.CreateTime == 0 {
		fileStore.BaseModel.CreateTime = util.CurMillionSeconds()
		fileStore.BaseModel.ModifyTime = util.CurMillionSeconds()
	}
	node, _ := snowflake.NewNode(0)
	id := node.Generate().String()
	fileStore.ID = id

	log.Printf("after fileStore proerties fill, fileStore is %+v \n", fileStore)
	config.DB.Create(&fileStore)
	return &fileStore, nil
}

func QueryFileRecord(queryEntity entity.FileStore) entity.FileStore {
	fileStore := &entity.FileStore{}
	config.DB.Where(&queryEntity).First(fileStore)
	return *fileStore
}
