package fileUtil

import (
	"context"
	"fmt"
	"log"
	"muses-engine/config"
	"muses-engine/internal/model/bo"
	"muses-engine/pkg/util"
	"time"

	"github.com/minio/minio-go/v7"
)

func UploadFile(minioUploadTask *bo.MinioUploadTask) error {
	bucketName := minioUploadTask.BucketName
	// file := minioUploadTask.FileInfo.File
	fileName := minioUploadTask.FileInfo.Name
	fileSize := minioUploadTask.FileInfo.Size

	log.Printf("upload file %s to minio, bucket name %s \n", fileName, minioUploadTask.BucketName)
	ctx := context.Background()
	contentType := "application/octet-stream"

	objectName := formatObjectName(minioUploadTask)
	minioUploadTask.ObjectName = objectName
	var info minio.UploadInfo
	var err error
	if minioUploadTask.LocalFileFlag {
		info, err = config.MinioClient.FPutObject(ctx, bucketName, objectName, minioUploadTask.FileInfo.FilePath, minio.PutObjectOptions{ContentType: contentType})
	} else {
		info, err = config.MinioClient.PutObject(ctx, bucketName, objectName, minioUploadTask.FileInfo.File, fileSize, minio.PutObjectOptions{ContentType: contentType})
	}
	if err != nil {
		log.Println("put object to minio failure ", err)
		return err
	}
	//info的返回参数里只有对象名，bucket是有用的，因为接下来要用它下载文件，etag是校验值类似MD5不能下载文件
	log.Println("Success upload, upload finish info is ", info)
	log.Printf("Successfully upload %s of size \t\t %dB \t %.2fKB \n ", objectName, info.Size, float64(info.Size/1024.0))
	return nil
}

// 下载可以考虑使用预签名URL生成可以访问的链接丢给前端，
func PreSignUrl(bucketName string, objectName string, expires time.Duration) (string, error) {
	url, err := config.MinioClient.PresignedGetObject(context.Background(), bucketName, objectName, expires, nil)
	if err != nil {
		log.Println("error while pre sign url ", err)
		return "", err
	}
	log.Println("pre sign url success ", url.String())
	return url.String(), nil
}

func DownloadFile(bucketName string, objectName string, filePath string) error {
	err := config.MinioClient.FGetObject(context.Background(), bucketName, objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("error while download file to local")
		return err
	}
	return nil
}

func createBucket(bucketName string) {
	location := "use-east-1"
	ctx := context.Background()

	err := config.MinioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})

	if err != nil {
		exists, errBucketExists := config.MinioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("we already own %s \n", bucketName)
		} else {
			if exists {
				log.Println("bucket is exists, err is ", err)
			} else {
				log.Fatalf("error under check bucket exists, err is %s errBucketExists %s \n", err, errBucketExists)
			}
		}
	} else {
		log.Printf("succcess create %s \n", bucketName)
	}

}

func formatObjectName(minioUploadTask *bo.MinioUploadTask) string {
	fileName := minioUploadTask.FileInfo.Name
	fileType := minioUploadTask.FileInfo.Type

	userId := minioUploadTask.UploadUserId
	uploadDir := config.MinioConfig.UploadDir
	log.Println("the config minio upload dir is ", uploadDir)
	filePathFormat := config.FileSystemConfig.UploadConfig.FilePathFormat
	fileNameFormat := config.FileSystemConfig.UploadConfig.FileNameFormat

	curDate := util.CurDate()
	minioUploadTask.TimeStamp = util.CurMillionSeconds()

	minioFileName := fmt.Sprintf(fileNameFormat, userId, minioUploadTask.TimeStamp, fileName)
	objectName := fmt.Sprintf(filePathFormat, uploadDir, fileType, curDate, minioFileName)
	log.Println("to format minio object name -> ", objectName)
	return objectName
}
