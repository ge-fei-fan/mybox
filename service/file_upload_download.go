package service

import (
	"errors"
	"gorm.io/gorm"
	"mime/multipart"
	"mybox/global"
	"mybox/model/system"
	"mybox/utils/upload"
	"path"
	"strconv"
)

type FileUploadAndDownloadService struct{}

func (fs *FileUploadAndDownloadService) Upload(file *system.FileUploadAndDownload) error {
	return global.BOX_DB.Create(file).Error
}

func (fs *FileUploadAndDownloadService) UploadFile(pathid string, fileHeader *multipart.FileHeader) (file *system.FileUploadAndDownload, err error) {
	oss := upload.NewOss()
	uploadUrl, uploadName, uploadErr := oss.UploadFile(fileHeader)
	if uploadErr != nil {
		panic(err)
	}
	id, _ := strconv.Atoi(pathid)
	f := &system.FileUploadAndDownload{
		SysUserRepositoryID: uint(id),
		Name:                fileHeader.Filename,
		Url:                 uploadUrl,
		Ext:                 path.Ext(fileHeader.Filename),
		Key:                 uploadName,
		Size:                fileHeader.Size,
	}
	return f, fs.Upload(f)
}

func (fs *FileUploadAndDownloadService) CheckDirId(pathId interface{}, userId uint) error {
	if global.BOX_DB == nil {
		return errors.New("db not init")
	}
	var sur system.SysUserRepository
	err := global.BOX_DB.Where("id=? ", pathId).First(&sur).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文件夹不存在")
		}
		return err
	}

	if sur.SysUserID != userId {
		return errors.New("没有文件夹权限")
	}
	return nil

}

func (fs *FileUploadAndDownloadService) List(userId, rpositoryId uint) (*[]system.SysUserRepository, *[]system.FileUploadAndDownload, error) {
	var rpositoryAll []system.SysUserRepository
	var files []system.FileUploadAndDownload
	err := global.BOX_DB.Where("sys_user_id=? and parent_id =?", userId, rpositoryId).Find(&rpositoryAll).Error
	if err != nil {
		return &rpositoryAll, &files, err
	}
	err = global.BOX_DB.Where("sys_user_repository_id=?", rpositoryId).Find(&files).Error
	if err != nil {
		return &rpositoryAll, &files, err
	}
	return &rpositoryAll, &files, nil
}
