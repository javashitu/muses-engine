package message

type VideoPubMsg struct {
	VideoProgramId string `json:"videoProgramId"`
	FileStoreId    string `json:"fileStoreId"`
	PubUserId      string `json:"pubUserId"`
}

// 有需要再加上转码状态，现在直接作为转码完毕的消息
type VideoTrscodeMsg struct {
	VideoProgramId string `json:"videoProgramId"`
	FileStoreId    string `json:"fileStoreId"`
	VideoMetaId    string `json:"videoMetaId"`
	PubUserId      string `json:"pubUserId"`
}
