package system

type FileUploadAndDownload struct {
	BOX_MODEL
	SysUserRepositoryID uint   `json:"sysUserRepositoryID" gorm:"comment:关联文件夹ID"`
	Name                string `json:"name" gorm:"comment:文件名"`  // 文件名
	Url                 string `json:"url" gorm:"comment:文件地址"`  // 文件地址
	Ext                 string `json:"ext" gorm:"comment:文件后缀"`  // 文件后缀
	Key                 string `json:"key" gorm:"comment:编号"`    // 编号，唯一标识
	Size                int64  `json:"size" gorm:"comment:文件大小"` //文件大小
}

type ChangeFile struct {
	Key  string `json:"key" gorm:"comment:编号"`   // 编号，唯一标识
	Name string `json:"name" gorm:"comment:文件名"` // 重命名的文件名
}

func (FileUploadAndDownload) TableName() string {
	return "file_upload_and_downloads"
}
