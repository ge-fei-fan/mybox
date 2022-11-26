package response

import (
	"github.com/golang-jwt/jwt/v4"
	"mybox/model/system"
)

type SysUserResponse struct {
	User system.SysUser `json:"user"`
}

type LoginResponse struct {
	User      system.SysUser   `json:"user"`
	Token     string           `json:"token"`
	ExpiresAt *jwt.NumericDate `json:"expiresAt"`
}
