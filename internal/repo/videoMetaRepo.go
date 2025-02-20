package repo

import (
	"log"
	"muses-engine/config"
	"muses-engine/internal/model/entity"
	"muses-engine/pkg/util"

	"github.com/bwmarrin/snowflake"
)

func SaveVideoMeta(videoMeta *entity.VideoMeta) (*entity.VideoMeta, error) {
	if videoMeta.CreateTime == 0 {
		videoMeta.CreateTime = util.CurMillionSeconds()
		videoMeta.ModifyTime = util.CurMillionSeconds()
	}
	node, _ := snowflake.NewNode(0)
	id := node.Generate().String()
	videoMeta.ID = id

	log.Printf("after videoMeta proerties fill, videoMeta is %+v \n", videoMeta)
	config.DB.Create(&videoMeta)
	return videoMeta, nil
}
