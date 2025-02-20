package entity

type FileStore struct {
	BaseModel
	ID         string `gorm:"column:id" json:"id"`
	FileName   string `gorm:"column:file_name"`
	Type       string `gorm:"column:type"`
	Size       int64  `gorm:"column:size"`
	UserId     string `gorm:"column:user_id"`
	BucketName string `gorm:"column:bucket_name"`
	ObjectName string `gorm:"column:object_name"`
}

func (FileStore) TableName() string {
	return "file_store"
}

func (this FileStore) IsNil() bool {
	return this == FileStore{}
}
