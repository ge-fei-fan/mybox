package service

import (
	"errors"
	"gorm.io/gorm"
	"mime/multipart"
	"mybox/global"
	"mybox/model/system"
	systemResp "mybox/model/system/response"
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

func (fs *FileUploadAndDownloadService) List(userId, rpositoryId uint) (*systemResp.ListItem, error) {
	var rpositoryAll []system.SysUserRepository
	var files []system.FileUploadAndDownload
	listItem := systemResp.ListItem{}
	if global.BOX_DB == nil {
		return &listItem, errors.New("db not init")
	}
	err := global.BOX_DB.Where("sys_user_id=? and parent_id =?", userId, rpositoryId).Find(&rpositoryAll).Error
	if err != nil {
		return &listItem, err
	}
	err = global.BOX_DB.Where("sys_user_repository_id=?", rpositoryId).Find(&files).Error
	if err != nil {
		return &listItem, err
	}
	p, err := fs.GetAllPath(rpositoryId)
	if err != nil {
		return &listItem, err
	}
	listItem.Path = p
	listItem.Folder = rpositoryAll
	listItem.File = files
	return &listItem, nil
}

func (fs *FileUploadAndDownloadService) GetAllPath(id uint) (fullPath string, err error) {
	sur := system.SysUserRepository{}
	p := ""
	err = global.BOX_DB.Where("id=?", id).First(&sur).Error
	if err == nil {
		if sur.ParentId != 0 {
			p, _ = fs.GetAllPath(sur.ParentId)
		}
		p += sur.DirName + "/"
		return p, nil
	}
	return "", err
}
