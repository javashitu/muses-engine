package entity

type BaseModel struct {
	CreateTime int64  `gorm:"column:create_time"`
	ModifyTime int64  `gorm:"column:modify_time"`
	DelFlag    string `gorm:"column:del_flag"`
}
