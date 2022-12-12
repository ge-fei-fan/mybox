package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mime/multipart"
	"mybox/global"
	"mybox/model/system"
	systemReq "mybox/model/system/request"
	systemResp "mybox/model/system/response"
	"mybox/utils"
	"mybox/utils/upload"
	"path"
	"strconv"
	"time"
)

type FileUploadAndDownloadService struct{}

func (fs *FileUploadAndDownloadService) Upload(file *system.FileUploadAndDownload) error {
	return global.BOX_DB.Create(file).Error
}

func (fs *FileUploadAndDownloadService) UploadFile(pathid string, fileHeader *multipart.FileHeader) (file *system.FileUploadAndDownload, err error) {
	err = fs.CheckFileName(fileHeader.Filename)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			oss := upload.NewOss()
			uploadUrl, uploadName, uploadErr := oss.UploadFile(fileHeader)
			if uploadErr != nil {
				panic(uploadErr)
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
		return nil, err
	}
	return nil, errors.New("存在相同文件名")

}
func (fs *FileUploadAndDownloadService) BreakPointContinue(chunkNum string, fileHeader *multipart.FileHeader) error {
	oss := upload.NewOss()
	err := oss.BreakPointContinue(chunkNum, fileHeader)
	if err != nil {
		return err
	}
	//global.BOX_REDIS.HSet(context.Background(), "chunk:"+fileHeader.Filename, chunkNum, "Done")
	err = global.BOX_REDIS.SAdd(context.Background(), "chunk:"+fileHeader.Filename, chunkNum).Err()
	return err
}
func (fs *FileUploadAndDownloadService) MergeFile(fileName string, FileMd5 string) (file *system.FileUploadAndDownload, err error) {
	var bfile system.File
	key := "file:" + fileName
	if err := global.BOX_REDIS.HGetAll(context.Background(), key).Scan(&bfile); err != nil {
		global.BOX_LOG.Error("MergeFile HGetAll() ", zap.Error(err))
		return nil, errors.New("MergeFile HGetAll(), err:" + err.Error())
	}
	oss := upload.NewOss()
	uploadUrl, uploadName, uploadErr := oss.MergeFile(fileName, FileMd5)
	if uploadErr != nil {
		return nil, uploadErr
	}
	id, _ := strconv.Atoi(bfile.FilePath)
	size, err := strconv.ParseInt(bfile.Size, 10, 64)
	f := &system.FileUploadAndDownload{
		SysUserRepositoryID: uint(id),
		Name:                fileName,
		Url:                 uploadUrl,
		Ext:                 path.Ext(fileName),
		Key:                 uploadName,
		Size:                size,
	}
	return f, fs.Upload(f)
}

func (fs *FileUploadAndDownloadService) DeleteFile(file *system.ChangeFile) (err error) {
	var fileDb system.FileUploadAndDownload
	err = global.BOX_DB.Where("`key` = ?", file.Key).First(&fileDb).Error
	if err != nil {
		return
	}
	oss := upload.NewOss()
	if err = oss.DeleteFile(fileDb.Key); err != nil {
		return errors.New("文件删除失败")
	}
	err = global.BOX_DB.Unscoped().Delete(&fileDb).Error
	return err
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

func (fs *FileUploadAndDownloadService) CheckFileName(name string) (err error) {
	return global.BOX_DB.Where("name=?", name).First(&system.FileUploadAndDownload{}).Error
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

func (fs *FileUploadAndDownloadService) RenameFile(key, name string) (file *system.FileUploadAndDownload, err error) {
	var sysfile system.FileUploadAndDownload
	if global.BOX_DB == nil {
		return nil, errors.New("db not init")
	}
	err = global.BOX_DB.Where("`key` = ?", key).First(&sysfile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文件不存在")
		}
		return nil, err
	}
	err = global.BOX_DB.Where("sys_user_repository_id=? and name=?", sysfile.SysUserRepositoryID, name).First(&system.FileUploadAndDownload{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = global.BOX_DB.Model(&sysfile).Update("name", name).Error
			return &sysfile, err
		}
		return &sysfile, err
	}
	return &sysfile, errors.New("当前目录下存在相同文件名")
}

func (fs *FileUploadAndDownloadService) ShareFile(sq *systemReq.ShareRequest) (file *system.SharePool, err error) {
	var sf system.SharePool
	sf.UUID = uuid.New()
	if sq.ExpiredTime != "" {
		exp, _ := utils.ParseDuration(sq.ExpiredTime)
		sf.ExpiredTime = exp
	}
	sf.FileKey = sq.File
	err = global.BOX_DB.Create(&sf).Error
	return &sf, err
}
func (fs *FileUploadAndDownloadService) GetShareFile(uuid string) (*system.FileUploadAndDownload, error) {
	var file system.FileUploadAndDownload
	var share system.SharePool
	filekey, err := fs.GetRedisShare(uuid)
	if err != nil {
		global.BOX_LOG.Error("获取分享文件查redis失败", zap.Error(err))
	}
	if filekey == "" {
		err := global.BOX_DB.Where("uuid=?", uuid).First(&share).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("分享链接不存在")
			}
			return nil, err
		}
		//判断一下是否过期
		if share.ExpiredTime != 0 {
			if time.Now().Sub(share.CreatedAt) > share.ExpiredTime {
				//已过期
				return nil, errors.New("分享链接已过期")
			}
			println(share.CreatedAt.Add(share.ExpiredTime).Unix())
			println(time.Now().Unix())
		}
		//没过期重新设置redis
		if err := fs.SetRedisShare(&share); err != nil {
			global.BOX_LOG.Error("设置分享redis失败", zap.Error(err))
		}
		filekey = share.FileKey
	}
	err = global.BOX_DB.Where("`key`=?", filekey).First(&file).Error
	return &file, err

}
func (fs *FileUploadAndDownloadService) SetRedisShare(share *system.SharePool) error {
	exp := share.ExpiredTime - time.Now().Sub(share.CreatedAt)
	return global.BOX_REDIS.Set(context.Background(), "share:"+share.UUID.String(), share.FileKey, exp).Err()
}

func (fs *FileUploadAndDownloadService) GetRedisShare(uuid string) (fileKey string, err error) {
	return global.BOX_REDIS.Get(context.Background(), "share:"+uuid).Result()
}
