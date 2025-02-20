package entity

import (
	"log"
	"muses-engine/internal/model/bo"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/vansante/go-ffprobe.v2"
)

type VideoMeta struct {
	BaseModel
	ID          string `gorm:"column:id" json:"id"`
	FileName    string `gorm:"column:file_name"`
	ProgramId   string `gorm:"column:program_id"`
	FileStoreId string `gorm:"column:file_store_id"`
	Duration    int    `gorm:"column:duration"`
	Format      string `gorm:"column:format"`
	Suffix      string `gorm:"column:suffix"`
	FrameRate   string `gorm:"column:frame_rate"`
	Resolution  string `gorm:"column:resolution"`
	Width       int    `gorm:"column:width"`
	Height      int    `gorm:"column:height"`
	BitRate     int64  `gorm:"column:bit_rate"`
	Size        int64  `gorm:"column:size"`
	VideoCodec  string `gorm:"column:video_codec"`
	AudioCodec  string `gorm:"column:audio_codec"`
}

func (VideoMeta) TableName() string {
	return "video_meta"
}

func GenVideoMeta(probeData *ffprobe.ProbeData, programId string, videoStoreId string) *VideoMeta {

	videoMeta := &VideoMeta{}

	size, err := strconv.ParseInt(probeData.Format.Size, 10, 64)
	if err != nil {
		log.Printf("resolve video's size error, the size is %s the error is %s", probeData.Format.Size, err)
	}
	videoMeta.Size = size
	videoMeta.Duration = int(probeData.Format.DurationSeconds)
	bitRate, err := strconv.ParseInt(probeData.Format.BitRate, 10, 64)
	if err != nil {
		log.Printf("resolve video's bitRate error, the bitRate is %s ,the error is %s ", probeData.Format.BitRate, err)
	}
	videoMeta.BitRate = bitRate
	videoMeta.FrameRate = probeData.Streams[0].RFrameRate

	videoMeta.VideoCodec = probeData.Streams[0].CodecName
	videoMeta.AudioCodec = probeData.Streams[1].CodecName

	// file, _ := os.Open(probeData.Format.Filename)
	// defer file.Close()

	videoMeta.FileName = filepath.Base(probeData.Format.Filename)
	videoMeta.Format = probeData.Format.FormatName
	videoMeta.Suffix = probeData.Format.Filename[strings.LastIndex(probeData.Format.Filename, ".")+1:]
	videoMeta.Width = probeData.Streams[0].Width
	videoMeta.Height = probeData.Streams[0].Height
	videoMeta.Resolution = bo.WrapResolution(videoMeta.Height).Name
	videoMeta.FileStoreId = videoStoreId
	videoMeta.ProgramId = programId

	return videoMeta

}
