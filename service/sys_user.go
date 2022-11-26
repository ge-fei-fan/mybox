package service

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mybox/global"
	"mybox/model/system"
	"mybox/utils"
)

type UserService struct{}

func (us *UserService) Register(su *system.SysUser) (userInter *system.SysUser, err error) {
	if global.BOX_DB == nil {
		return nil, errors.New("db not init")
	}

	var user system.SysUser
	err = global.BOX_DB.Where("username=?", su.Username).First(&user).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		//用户已注册或者查询出错
		return nil, errors.New("用户名已注册")
	}
	//创建一个用户名同名的文件夹
	sur := system.SysUserRepository{
		DirName: su.Username,
	}
	su.Repository = []system.SysUserRepository{sur}
	su.UUID = uuid.New()
	su.Password = utils.BcryptHash(su.Password)
	err = global.BOX_DB.Create(&su).Error
	return su, err
}

func (us *UserService) Login(su *system.SysUser) (userInter *system.SysUser, err error) {
	if global.BOX_DB == nil {
		return nil, errors.New("db not init")
	}
	var user system.SysUser
	err = global.BOX_DB.Where("username =?", su.Username).First(&user).Error
	if err == nil {
		if !utils.BcryptCheck(su.Password, user.Password) {
			return nil, errors.New("密码错误")
		}
	}
	return &user, err
}

func (us *UserService) GetUserInfo(uuid uuid.UUID) (user system.SysUser, err error) {
	var reqUser system.SysUser
	err = global.BOX_DB.First(&reqUser, "uuid = ?", uuid).Error
	if err != nil {
		return reqUser, err
	}
	return reqUser, err
}
