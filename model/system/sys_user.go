package system

import (
	"github.com/google/uuid"
)

type SysUser struct {
	BOX_MODEL
	UUID       uuid.UUID           `json:"uuid" gorm:"index;comment:用户UUID"`                                                     // 用户UUID
	Username   string              `json:"userName" gorm:"index;comment:用户登录名"`                                                  // 用户登录名
	Password   string              `json:"-"  gorm:"comment:用户登录密码"`                                                             // 用户登录密码
	NickName   string              `json:"nickName" gorm:"default:系统用户;comment:用户昵称"`                                            // 用户昵称
	HeaderImg  string              `json:"headerImg" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"` // 用户头像
	Email      string              `json:"email"  gorm:"comment:用户邮箱"`                                                           // 用户邮箱
	Repository []SysUserRepository `json:"repository" gorm:"comment:用户关联的文件夹和文件"`
}

func (su *SysUser) TableName() string {
	return "sys_users"
}
