package v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
		global.BOX_LOG.Error("修改数据库链接失败", zap.Error(err))
		response.FailWithMessage("修改数据库链接失败", c)
		return
	}
	response.OkWithDetailed(systemResp.FileResponse{File: *file}, "上传成功", c)
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
		global.BOX_LOG.Error("查询出错", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	}

	response.OkWithDetailed(list, "查询成功", c)
}
