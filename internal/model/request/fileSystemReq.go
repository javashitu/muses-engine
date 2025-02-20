package request

type UploadFileReq struct {
	UserId   string `form:"userId" binding:"required"`
	UserName string `form:"userName"`
	// BucketName string
	// //下面的字段是返回结果用
	// FileName   string
	// Type       string
	// Size       int64
	// ObjectName string
	// TimeStamp  int64
}

type VisitFileUrlReq struct {
	Id            string `json:"id"`
	ExpireSeconds int    `json:"expireSeconds"`
}
