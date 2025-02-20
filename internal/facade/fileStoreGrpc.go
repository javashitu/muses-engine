package facade

import (
	"context"
	"log"
	"muses-engine/config"
	"muses-engine/internal/fileUtil"
	"muses-engine/internal/model/entity"
	"muses-engine/internal/repo"
	"muses-engine/pkg/proto_generated/pb"
	"time"
)

type MyHello struct {
	pb.UnimplementedMyHelloServer
}

func (this *MyHello) mustEmbedUnimplementedHelloServer() {
	//TODO implement me
	panic("implement me")
}

func (this *MyHello) SayHello(ctx context.Context, p *pb.Person) (*pb.Person, error) {
	log.Println(p.Name)
	p.Name = "hello" + p.Name
	return p, nil
}

type FileStoreService struct {
	pb.UnimplementedFileStoreServiceServer
}

func (fileStoreService *FileStoreService) mustEmbedUnimplementedFileStoreService() {
	panic("need implements FileStoreService")
}

func (fileStoreService *FileStoreService) QueryVisitUrl(ctx context.Context, request *pb.QueryFileInfoReq) (*pb.QueryFileInfoRsp, error) {
	log.Println("begin query file visiturl, the query request ", request)
	fileStoreInfoSlice := make([]*pb.FileStoreInfo, len(request.QueryFileInfoList))
	for index, value := range request.QueryFileInfoList {
		queryEntity := entity.FileStore{}
		queryEntity.ID = value.Id
		fileStore := repo.QueryFileRecord(queryEntity)

		bucketName := fileStore.BucketName
		objectName := fileStore.ObjectName
		expireSeconds := int(value.PreviewExpireSeconds)
		if expireSeconds == 0 || expireSeconds > config.MinioConfig.AccessExpires {
			expireSeconds = config.MinioConfig.AccessExpires
		}

		url, err := fileUtil.PreSignUrl(bucketName, objectName, time.Duration(expireSeconds)*time.Second)
		if err != nil {
			log.Printf("generate visiturl of %s failure \r\n ", value.Id)
			continue
		}
		fileStoreInfo := pb.FileStoreInfo{}
		fileStoreInfo.Id = value.Id
		fileStoreInfo.Url = url

		fileStoreInfoSlice[index] = &fileStoreInfo
	}
	queryFileRsp := &pb.QueryFileInfoRsp{}
	queryFileRsp.FileStoreInfoList = fileStoreInfoSlice
	return queryFileRsp, nil

}
