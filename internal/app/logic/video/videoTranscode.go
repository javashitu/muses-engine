package videoTranscode

import (
	"context"
	"encoding/json"
	"log"
	"muses-engine/config"
	"muses-engine/internal/app/logic/fileSystem"
	"muses-engine/internal/fileUtil"
	"muses-engine/internal/kafka/producer"
	"muses-engine/internal/model/entity"
	kafkaMsg "muses-engine/internal/model/message"
	"muses-engine/internal/repo"
	"os/exec"
	"time"

	"gopkg.in/vansante/go-ffprobe.v2"
)

func TranscodeVideoFile(videoPubMsg kafkaMsg.VideoPubMsg) {
	// videoFileId := videoPubMsg.VideoStoreId
	userId := videoPubMsg.PubUserId
	queryEntity := entity.FileStore{}
	queryEntity.ID = videoPubMsg.FileStoreId
	fileStore := repo.QueryFileRecord(queryEntity)
	downloadDir := config.FileSystemConfig.DownloadConfig.VideoDownloadDir
	filePath := downloadDir + "/" + fileStore.FileName

	err := fileUtil.DownloadFile(fileStore.BucketName, fileStore.ObjectName, filePath)
	if err != nil {
		log.Println("download video file failure, download path ", filePath)
		return
	}
	outputFilePath := downloadDir + "/360p" + fileStore.FileName
	err = transcodeVideo(filePath, outputFilePath)
	if err != nil {
		log.Println("transcode video failure, error is ", err)
		// return
	}
	savedFileStore := fileSystem.SaveLocalFile(outputFilePath, userId)

	if savedFileStore == nil {
		log.Println("can't find localfile and store, so won't resolve, the outputFilePath is ", outputFilePath)
	}
	videoMeta := resolveVideo(outputFilePath, videoPubMsg.VideoProgramId, savedFileStore.ID)
	if videoMeta == nil {
		log.Println("resolve video failure ,the resolve filePath is ", outputFilePath)
		return
	}
	savedVideoMeta, _ := repo.SaveVideoMeta(videoMeta)

	videoTrscodeMsg, _ := json.Marshal(genVideoTranscodeMsg(videoPubMsg, savedVideoMeta.ID))
	producer.SendTrsCodeFinishMsg(videoTrscodeMsg, "")
}

func genVideoTranscodeMsg(videoPubMsg kafkaMsg.VideoPubMsg, videoMetaId string) kafkaMsg.VideoTrscodeMsg {
	videoTrscodeMsg := kafkaMsg.VideoTrscodeMsg{}
	videoTrscodeMsg.PubUserId = videoPubMsg.PubUserId
	videoTrscodeMsg.VideoProgramId = videoPubMsg.VideoProgramId
	videoTrscodeMsg.FileStoreId = videoPubMsg.FileStoreId
	videoTrscodeMsg.VideoMetaId = videoMetaId
	return videoTrscodeMsg
}

func transcodeVideo(inputFilePath string, outputFilePath string) error {
	command := config.FileSystemConfig.VideoManager.FmtTrscodeVideoCmd(inputFilePath, outputFilePath)
	log.Println("will execute command to transcode video, execute command is ", command)
	cmd := exec.Command("sh", command)
	output, err := cmd.Output()
	log.Println("execute transcode video shell output ", output)

	if err != nil {
		log.Println("execute transcode video shell failure error is ", err)
		return err
	}
	return nil
}

func resolveVideo(filePath string, programId string, fileStoreId string) *entity.VideoMeta {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()
	//https://pkg.go.dev/gopkg.in/vansante/go-ffprobe.v2#ProbeData

	probeData, err := ffprobe.ProbeURL(ctx, filePath)
	if err != nil {
		log.Printf("Error resolve video data by ffprobe the error is %v", err)
		return nil
	}
	return entity.GenVideoMeta(probeData, programId, fileStoreId)
}
