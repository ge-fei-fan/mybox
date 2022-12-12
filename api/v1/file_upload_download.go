package v1

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mybox/global"
	"mybox/model/common/response"
	"mybox/model/system"
	systemReq "mybox/model/system/request"
	systemResp "mybox/model/system/response"
	"mybox/utils"
)

type FileUploadAndDownloadApi struct{}

func (f *FileUploadAndDownloadApi) UploadFile(c *gin.Context) {
	var file *system.FileUploadAndDownload
	pathId, ok := c.GetPostForm("path")
	if !ok {
		response.FailWithMessage("path参数为空", c)
		return
	}
	claimsId := utils.GetUserId(c)
	err := fileService.CheckDirId(pathId, claimsId)
	if err != nil {
		global.BOX_LOG.Error("文件夹不存在", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	_, header, err := c.Request.FormFile("file")
	if err != nil {
		global.BOX_LOG.Error("接收文件失败", zap.Error(err))
		response.FailWithMessage("接收文件失败", c)
		return
	}
	file, err = fileService.UploadFile(pathId, header)
	if err != nil {
		global.BOX_LOG.Error("上传文件失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(systemResp.FileResponse{File: *file}, "上传成功", c)
}

func (f *FileUploadAndDownloadApi) DeleteFile(c *gin.Context) {
	var file system.ChangeFile
	err := c.ShouldBindJSON(&file)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = fileService.DeleteFile(&file)
	if err != nil {
		global.BOX_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (f *FileUploadAndDownloadApi) List(c *gin.Context) {
	var cf systemReq.ChangeFolder
	err := c.ShouldBind(&cf)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if cf.ID == 0 {
		response.FailWithMessage("id必传", c)
		return
	}

	claimsId := utils.GetUserId(c)
	list, err := fileService.List(claimsId, cf.ID)
	if err != nil {
		global.BOX_LOG.Error("查询文件信息出错", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	}

	response.OkWithDetailed(list, "查询成功", c)
}

func (f *FileUploadAndDownloadApi) RenameFile(c *gin.Context) {
	var cf system.ChangeFile
	err := c.ShouldBind(&cf)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if cf.Key == "" || cf.Name == "" {
		response.FailWithMessage("参数不足", c)
		return
	}
	file, err := fileService.RenameFile(cf.Key, cf.Name)
	if err != nil {
		global.BOX_LOG.Error("修改文件名出错", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	}
	response.OkWithDetailed(*file, "修改成功", c)
}

func (f *FileUploadAndDownloadApi) ShareFile(c *gin.Context) {
	var sq systemReq.ShareRequest
	err := c.ShouldBind(&sq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if sq.File == "" {
		response.FailWithMessage("文件为空", c)
		return
	}
	share, err := fileService.ShareFile(&sq)
	if err != nil {
		global.BOX_LOG.Error("创建分享数据失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	}
	if err := fileService.SetRedisShare(share); err != nil {
		global.BOX_LOG.Error("设置分享redis失败", zap.Error(err))
	}
	response.FailWithDetailed(*share, "创建分享成功", c)
}
func (f *FileUploadAndDownloadApi) ShareFileInfo(c *gin.Context) {
	var si systemReq.ShareidRequest
	err := c.ShouldBind(&si)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if si.ShareId == "" {
		response.FailWithMessage("shareid为空", c)
		return
	}

	fileinfo, err := fileService.GetShareFile(si.ShareId)
	if err != nil {
		global.BOX_LOG.Error("获取分享文件信息失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(*fileinfo, "获取分享信息成功", c)
}

func (f FileUploadAndDownloadApi) Check(c *gin.Context) {
	var redisFile system.File
	err := c.ShouldBind(&redisFile)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if redisFile.FileName == "" || redisFile.FilePath == "" || redisFile.FileMd5 == "" {
		response.FailWithMessage("缺少必填参数", c)
		return
	}
	err = fileService.CheckFileName(redisFile.FileName)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			response.FailWithMessage(err.Error(), c)
			return
		}
	} else {
		response.FailWithMessage("存在相同文件名", c)
		return
	}
	_, err = global.BOX_REDIS.Pipelined(context.Background(), func(rdb redis.Pipeliner) error {
		rdb.HSet(context.Background(), "file:"+redisFile.FileName, "FileMd5", redisFile.FileMd5)
		rdb.HSet(context.Background(), "file:"+redisFile.FileName, "FilePath", redisFile.FilePath)
		rdb.HSet(context.Background(), "file:"+redisFile.FileName, "ChunkTotal", redisFile.ChunkTotal)
		rdb.HSet(context.Background(), "file:"+redisFile.FileName, "Size", redisFile.Size)
		return nil
	})
	if err != nil {
		global.BOX_LOG.Error("分片上传设置redis文件信息失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("文件详情获取成功", c)
}

func (f FileUploadAndDownloadApi) BreakPointContinue(c *gin.Context) {
	chunkNumber := c.Request.FormValue("chunkNumber")
	_, FileHeader, err := c.Request.FormFile("file")
	if err != nil {
		global.BOX_LOG.Error("接收文件失败!", zap.Error(err))
		response.FailWithMessage("接收文件失败", c)
		return
	}
	err = fileService.BreakPointContinue(chunkNumber, FileHeader)
	if err != nil {
		global.BOX_LOG.Error("保存切片出错", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("保存切片成功", c)
}
func (f FileUploadAndDownloadApi) BreakPointContinueFinish(c *gin.Context) {
	fileMd5 := c.Request.FormValue("fileMd5")
	fileName := c.Request.FormValue("fileName")
	if fileMd5 == "" || fileName == "" {
		response.FailWithMessage("缺少必填参数", c)
		return
	}
	file, err := fileService.MergeFile(fileName, fileMd5)
	if err != nil {
		global.BOX_LOG.Error("分片上传文件失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(systemResp.FileResponse{File: *file}, "上传成功", c)

}
