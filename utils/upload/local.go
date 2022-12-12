package upload

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"mybox/global"
	"mybox/model/system"
	"mybox/utils"
	"os"
	"path"
	"strconv"
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
	p := path.Join(global.BOX_CONFIG.Local.StorePath, key)
	if err := os.Remove(p); err != nil {
		return errors.New("本地文件删除失败, err:" + err.Error())
	}
	return nil
}

func (l *Local) BreakPointContinue(chunkNum string, file *multipart.FileHeader) error {

	newfilename := file.Filename + "_" + chunkNum
	mkdirErr := os.MkdirAll(global.BOX_CONFIG.Local.BreakPointPath, os.ModePerm)
	if mkdirErr != nil {
		global.BOX_LOG.Error("BreakPointContinue os.MkdirAll() Filed", zap.Error(mkdirErr))
		return errors.New("BreakPointContinue os.MkdirAll() Filed, err:" + mkdirErr.Error())
	}
	breakpointPath := path.Join(global.BOX_CONFIG.Local.BreakPointPath, newfilename)

	f, openError := file.Open() // 读取文件
	if openError != nil {
		global.BOX_LOG.Error("BreakPointContinue file.Open() Filed", zap.Error(openError))
		return errors.New("BreakPointContinue file.Open() Filed, err:" + openError.Error())
	}
	defer f.Close()

	out, createErr := os.Create(breakpointPath)
	if createErr != nil {
		global.BOX_LOG.Error("BreakPointContinue file.Create() Filed", zap.Error(createErr))
		return errors.New("BreakPointContinue file.Create() Filed, err:" + createErr.Error())
	}
	defer out.Close()

	//拷贝文件
	_, copyErr := io.Copy(out, f)
	if copyErr != nil {
		global.BOX_LOG.Error("BreakPointContinue file.Copy() Filed", zap.Error(copyErr))
		return errors.New("BreakPointContinue file.Copy() Filed, err:" + copyErr.Error())
	}
	return nil
}
func (l *Local) MergeFile(fileName string, FileMd5 string) (string, string, error) {
	var file system.File
	if err := global.BOX_REDIS.HGetAll(context.Background(), "file:"+fileName).Scan(&file); err != nil {
		global.BOX_LOG.Error("MergeFile HGetAll() ", zap.Error(err))
		return "", "", errors.New("MergeFile HGetAll(), err:" + err.Error())
	}
	chunRedis, err := global.BOX_REDIS.SCard(context.Background(), "chunk:"+fileName).Result()
	if err != nil {
		global.BOX_LOG.Error("MergeFile SCard() ", zap.Error(err))
		return "", "", errors.New("MergeFile SCard(), err:" + err.Error())
	}
	fmt.Println(chunRedis)
	chunRedisAll := global.BOX_REDIS.SMembers(context.Background(), "chunk:"+fileName).Val()

	fileChunkAll, _ := strconv.ParseInt(file.ChunkTotal, 10, 64)
	if chunRedis != fileChunkAll {
		global.BOX_LOG.Error("需要切片:"+file.ChunkTotal+"redis切片："+string(chunRedis), zap.Error(err))
		return "", "", errors.New("切片数量不足,请重新上传")
	}
	//读取文件后缀
	ext := path.Ext(fileName)
	//读取文件名并加密
	name := fileName
	name = utils.MD5V([]byte(name))
	//拼接新的文件名称
	newfilename := name + "_" + time.Now().Format("20060102150405") + ext

	_ = os.MkdirAll(global.BOX_CONFIG.Local.StorePath, os.ModePerm)
	//文件存储路径
	storePath := path.Join(global.BOX_CONFIG.Local.StorePath, newfilename)
	//访问路径
	filepath := path.Join(global.BOX_CONFIG.Local.Path, newfilename)
	//创建目标文件
	fd, err := os.OpenFile(storePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return "", "", err
	}
	defer fd.Close()

	//循环读取切片文件
	for _, chunknum := range chunRedisAll {
		content, _ := os.ReadFile(path.Join(global.BOX_CONFIG.Local.BreakPointPath, fileName+"_"+chunknum))
		_, err = fd.Write(content)
		if err != nil {
			//合并文件失败，把目标文件删了
			_ = os.Remove(storePath)
			return "", "", err
		}
		//合并成功把切片文件删了
		_ = os.Remove(path.Join(global.BOX_CONFIG.Local.BreakPointPath, fileName+"_"+chunknum))
	}
	//redis数据删了
	global.BOX_REDIS.Del(context.Background(), "file:"+fileName, "chunk:"+fileName)
	//比较哈希值
	sumMd5, err := utils.MD5(storePath)
	if err != nil {
		return "", "", errors.New("计算文件md5错误")
	}
	if FileMd5 != sumMd5 {
		return "", "", errors.New("合并文件后md5错误")
	}
	return filepath, newfilename, nil
}
