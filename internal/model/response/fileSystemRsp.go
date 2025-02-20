package response

type UploadFileRsp struct {
	ID string `json:"id"`
}

type VisitFiletUrlRsp struct {
	Url string `json:"url"`
}

type AccessUrlInfo struct {
	BucketName    string `json:"bucketName"`
	ObjectName    string `json:"objectName"`
	ExpireSeconds int32  `json:"expireSeconds"`
	Url           string `json:"url"`
}
