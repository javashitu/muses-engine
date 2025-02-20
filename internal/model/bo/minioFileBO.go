package bo

import (
	"io"
)

type MinioUploadTask struct {
	BucketName string
	//在上传之前这个值是没有拼接出来的
	ObjectName string

	// ObjectName   string
	FileInfo      FileInfo
	UploadUserId  string
	TimeStamp     int64
	LocalFileFlag bool
}

type MinioUplaodResult struct {
	BucketName string
	//下面的字段是返回结果用
	FileInfo FileInfo
	//上传时间
	TimeStamp    int64
	UploadUserId string
}

type FileInfo struct {
	// File     multipart.File
	File     io.Reader
	FilePath string
	Name     string
	Type     string
	Size     int64
}

func (this FileInfo) IsNil() bool {
	return this == FileInfo{}
}
