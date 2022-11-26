package system

type SysUserRepository struct {
	BOX_MODEL
	DirName   string                  `json:"dirName" gorm:"文件夹名称"`
	SysUserID uint                    `json:"sysUserID" gorm:"comment:关联的用户ID"`
	ParentId  uint                    `json:"parentId" gorm:"comment:父级目录ID 空为主目录"`
	Files     []FileUploadAndDownload `json:"files" gorm:"当前目录下保存的文件"`
}

func (sur *SysUserRepository) TableName() string {
	return "sys_user_repository"
}
