package upload

import (
	"errors"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"mybox/global"
	"mybox/utils"
	"os"
	"path"
	"time"
)

type Local struct{}

func (l *Local) UploadFile(file *multipart.FileHeader) (string, string, error) {
	//读取文件后缀
	ext := path.Ext(file.Filename)
	//读取文件名并加密
	filename := file.Filename
	filename = utils.MD5V([]byte(filename))
	//拼接新的文件名称
	newfilename := filename + "_" + time.Now().Format("20060102150405") + ext

	mkdirErr := os.MkdirAll(global.BOX_CONFIG.Local.StorePath, os.ModePerm)
	if mkdirErr != nil {
		global.BOX_LOG.Error("function os.MkdirAll() Filed", zap.Error(mkdirErr))
		return "", "", errors.New("function os.MkdirAll() Filed, err:" + mkdirErr.Error())
	}

	//拼接路径的文件名
	//文件存储路径
	storePath := path.Join(global.BOX_CONFIG.Local.StorePath, newfilename)
	//访问路径
	filepath := path.Join(global.BOX_CONFIG.Local.Path, newfilename)

	f, openError := file.Open() // 读取文件
	if openError != nil {
		global.BOX_LOG.Error("function file.Open() Filed", zap.Error(openError))
		return "", "", errors.New("function file.Open() Filed, err:" + openError.Error())
	}
	defer f.Close()

	out, createErr := os.Create(storePath)
	if createErr != nil {
		global.BOX_LOG.Error("function file.Create() Filed", zap.Error(createErr))
		return "", "", errors.New("function file.Create() Filed, err:" + createErr.Error())
	}
	defer out.Close()

	//拷贝文件
	_, copyErr := io.Copy(out, f)
	if copyErr != nil {
		global.BOX_LOG.Error("function file.Copy() Filed", zap.Error(copyErr))
		return "", "", errors.New("function file.Copy() Filed, err:" + copyErr.Error())
	}
	return filepath, newfilename, nil
}

func (l *Local) DeleteFile(key string) error {
	return nil
}
