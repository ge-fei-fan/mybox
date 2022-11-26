package service

import (
	"errors"
	"gorm.io/gorm"
	"mybox/global"
	"mybox/model/system"
)

type FolderService struct{}

func (fs *FolderService) CreateFolder(repository *system.SysUserRepository) (repositoryInter *system.SysUserRepository, err error) {
	if global.BOX_DB == nil {
		return repositoryInter, errors.New("db not init")
	}
	var sur system.SysUserRepository
	err = global.BOX_DB.Where("dir_name = ? AND parent_id >= ?", repository.DirName, repository.ID).First(&sur).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = global.BOX_DB.Create(repository).Error
			return repository, err
		}
		return repositoryInter, errors.New("创建文件夹出错")
	}
	return &sur, errors.New("文件夹名称已存在")
}

func (fs *FolderService) ChangeFolderName(id uint, dirname string) error {
	if global.BOX_DB == nil {
		return errors.New("db not init")
	}
	var sur system.SysUserRepository
	err := global.BOX_DB.Where("id = ?", id).First(&sur).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文件夹不存在")
		}
		return errors.New("操作出错")
	}
	err = global.BOX_DB.Where("parent_id = ? and dir_name =?", sur.ParentId, dirname).First(&system.SysUserRepository{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = global.BOX_DB.Model(&sur).Update("dir_name", dirname).Error
			return err
		}
		return err
	}

	return errors.New("该目录下文件夹名已存在")

}
