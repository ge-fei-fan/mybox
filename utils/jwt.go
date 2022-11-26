package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"mybox/global"
	"mybox/model/system/request"
	"time"
)

type MyJwt struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

func NewJwt() *MyJwt {
	return &MyJwt{
		SigningKey: []byte(global.BOX_CONFIG.JWT.SigningKey),
	}
}

func (mj *MyJwt) CreateClaims(baseClaims *request.BaseClaims) *request.CustomClaims {
	buftime, _ := ParseDuration(global.BOX_CONFIG.JWT.BufferTime)
	exptime, _ := ParseDuration(global.BOX_CONFIG.JWT.ExpiresTime)
	claims := &request.CustomClaims{
		BaseClaims: *baseClaims,
		BufferTime: int64(buftime / time.Second),
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exptime)),
			Issuer:    global.BOX_CONFIG.JWT.Issuer,
		},
	}

	return claims
}

func (mj *MyJwt) CreateToken(claims *request.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, *claims)
	return token.SignedString(mj.SigningKey)
}

func (mj *MyJwt) ParseToken(tokenString string) (*request.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &request.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mj.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors == jwt.ValidationErrorMalformed {
				return nil, TokenMalformed
			} else if ve.Errors == jwt.ValidationErrorExpired {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors == jwt.ValidationErrorNotValidYet {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*request.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}
}
