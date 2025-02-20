package config

import (
	"embed"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gopkg.in/yaml.v3"
)

var MinioConfig MinIO

type MinIO struct {
	EndPoint        string `yaml:"endPoint"`
	AccessKeyId     string `yaml:"accessKeyId"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	// LocalFile       string `yaml:"local_file"`
	UploadDir       string `yaml:"uploadDir"`
	UploadDirFormat string `yaml:"uploadDirFormat"`
	BucketName      string `yaml:"bucketName"`
	AccessExpires   int    `yaml:"accessExpires"`
}

var MinioClient *minio.Client

//go:embed minioConfig.yaml
var miof embed.FS

func init() {
	configMinio()
	initMinioClient()
}

func configMinio() {
	file, err := miof.ReadFile("minioConfig.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(file, &MinioConfig)
	if err != nil {
		panic(err)
	}
}

func initMinioClient() {
	endpoint := MinioConfig.EndPoint
	accessKeyID := MinioConfig.AccessKeyId
	secretAccessKey := MinioConfig.SecretAccessKey
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Panicf("connect to minio server failure,maybe server not run , err  %s ", err)
	}
	MinioClient = minioClient
}
