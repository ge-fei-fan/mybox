package system

import (
	"github.com/google/uuid"
	"time"
)

type SharePool struct {
	BOX_MODEL
	UUID        uuid.UUID     `json:"uuid" gorm:"index;comment:分享唯一UUID"`
	ExpiredTime time.Duration `json:"expiredTime" gorm:"comment:过期时间"` //为空永久
	CheckNum    int           `json:"checkNum" gorm:"comment:点击次数"`
	FileKey     string        `json:"fileKey" gorm:"comment:分享文件的唯一标识"`
}

func (SharePool) TableName() string {
	return "share_pool"
}
