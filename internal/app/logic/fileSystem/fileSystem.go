package fileSystem

import (
	"log"
	"mime/multipart"
	"muses-engine/config"
	"muses-engine/internal/common"
	"muses-engine/internal/fileUtil"
	"muses-engine/internal/model/bo"
	"muses-engine/internal/model/entity"
	"muses-engine/internal/model/request"
	"muses-engine/internal/model/response"
	"muses-engine/internal/repo"
	"os"
	"strings"
	"time"
)

func SaveFile(files []*multipart.FileHeader, uploadInfoRes *request.UploadFileReq) response.ApiResult {
	minioUploadTask := genMinioUploadTask(uploadInfoRes.UserId)
	minioUploadTask.FileInfo = genFileInfoFromNet(files[0])

	log.Printf("begin upload file, the upload task +%v \n ", minioUploadTask)
	err := fileUtil.UploadFile(minioUploadTask)
	if err != nil {
		return *response.FailureApi(common.UploadFileError)
	}

	fileStroe := genFileStroe(minioUploadTask, uploadInfoRes.UserId)
	saveFileStore, err := repo.SaveFileRecord(fileStroe)
	if err != nil {
		return *response.FailureApi(common.UploadFileError)
	}

	log.Println("the file store record -> ", saveFileStore)

	uploadFileRsp := response.UploadFileRsp{}
	uploadFileRsp.ID = saveFileStore.ID
	return *response.SuccessApi(uploadFileRsp)
}

func GenFileVisitUrl(visitFileUrlReq request.VisitFileUrlReq) response.ApiResult {
	log.Printf("begin generate file visit url +%v \n ", visitFileUrlReq)
	queryEntity := entity.FileStore{}
	queryEntity.ID = visitFileUrlReq.Id
	fileStore := repo.QueryFileRecord(queryEntity)

	bucketName := fileStore.BucketName
	objectName := fileStore.ObjectName
	expireSeconds := config.MinioConfig.AccessExpires
	if expireSeconds == 0 {
		expireSeconds = visitFileUrlReq.ExpireSeconds
	}

	url, err := fileUtil.PreSignUrl(bucketName, objectName, time.Duration(expireSeconds)*time.Second)
	if err != nil {
		return *response.FailureApi(common.PreSignUrlError)
	}
	queryFileVisitUrl := response.VisitFiletUrlRsp{}
	queryFileVisitUrl.Url = url

	return *response.SuccessApi(queryFileVisitUrl)
}

func SaveLocalFile(filePath string, userId string) *entity.FileStore {
	log.Printf("upload local file to minio, local file path %s ,upload userId %s \r\n", filePath, userId)
	minioUploadTask := genMinioUploadTask(userId)
	minioUploadTask.LocalFileFlag = true
	minioUploadTask.FileInfo = genFileInfoFromLocal(filePath)

	log.Printf("uplaod minio task is %v \r\n", minioUploadTask)
	if minioUploadTask.FileInfo.IsNil() {
		log.Println("no fileInfo in minioUploadTask, maybe gen fileInfo failure")
		return nil
	}

	err := fileUtil.UploadFile(minioUploadTask)
	if err != nil {
		log.Println("upload local file to minio failure, the error is ", err)
		return nil
	}
	fileStroe := genFileStroe(minioUploadTask, userId)

	saveFileStore, err := repo.SaveFileRecord(fileStroe)
	if err != nil {
		log.Println("save file to fileStore repo failure,the error is ", err)
		return nil
	}

	log.Println("the file store record -> ", saveFileStore)
	return saveFileStore
}

//=================================================================下方私有方法=================================================================================

func genMinioUploadTask(userId string) *bo.MinioUploadTask {
	minioUploadTask := &bo.MinioUploadTask{}
	minioUploadTask.BucketName = config.MinioConfig.BucketName
	minioUploadTask.UploadUserId = userId

	return minioUploadTask
}

func genFileInfoFromNet(file *multipart.FileHeader) bo.FileInfo {
	fileInfo := bo.FileInfo{}
	fileInfo.File, _ = file.Open()
	fileInfo.Name = file.Filename
	fileInfo.Size = file.Size
	index := strings.Index(file.Filename, ".")
	fileInfo.Type = file.Filename[index+1:]

	return fileInfo
}

func genFileInfoFromLocal(filePath string) bo.FileInfo {
	defer func() {
		if p := recover(); p != nil {
			log.Println("catch panic, the error is ", p)
		}
	}()
	file, err := os.Stat(filePath)
	if err != nil {
		log.Printf("open file failure, can't open file %s ", err)
	}
	//这两个file的类型完全不一样不能提取出公共的方法
	fileInfo := bo.FileInfo{}
	// fileInfo.File = file
	fileInfo.Name = file.Name()
	fileInfo.Size = file.Size()
	fileInfo.FilePath = filePath
	index := strings.Index(fileInfo.Name, ".")
	fileInfo.Type = fileInfo.Name[index+1:]
	return fileInfo
}

func genFileStroe(minioUploadTask *bo.MinioUploadTask, userId string) entity.FileStore {
	fileStore := entity.FileStore{}
	fileStore.BucketName = minioUploadTask.BucketName
	fileStore.ObjectName = minioUploadTask.ObjectName
	fileStore.FileName = minioUploadTask.FileInfo.Name
	fileStore.Size = minioUploadTask.FileInfo.Size
	fileStore.Type = minioUploadTask.FileInfo.Type
	fileStore.UserId = userId

	return fileStore
}
