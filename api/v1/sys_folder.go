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

type FolderApi struct{}

func (f *FolderApi) CreateFolder(c *gin.Context) {
	var cf systemReq.CreateFolder
	err := c.ShouldBind(&cf)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if len(cf.FolderName) == 0 || cf.ParentId == 0 {
		response.FailWithMessage("参数不足", c)
		return
	}
	claimsId := utils.GetUserId(c)
	err = fileService.CheckDirId(cf.ParentId, claimsId)
	if err != nil {
		global.BOX_LOG.Error("文件夹不存在", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	sur := &system.SysUserRepository{DirName: cf.FolderName, SysUserID: claimsId, ParentId: cf.ParentId}
	repositoryReturn, err := folderService.CreateFolder(sur)
	if err != nil {
		global.BOX_LOG.Error("创建文件夹失败", zap.Error(err))
		response.FailWithDetailed(systemResp.RepositoryResponse{Repository: *repositoryReturn}, "创建文件夹失败", c)
		return
	}
	response.OkWithDetailed(systemResp.RepositoryResponse{Repository: *repositoryReturn}, "创建文件夹成功", c)
}

func (f *FolderApi) ChangeFolderName(c *gin.Context) {
	var cf systemReq.ChangeFolder
	err := c.ShouldBind(&cf)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if cf.FolderName == "" || cf.ID == 0 {
		response.FailWithMessage("参数不足", c)
		return
	}
	//判断是否有文件夹权限
	claimsId := utils.GetUserId(c)
	err = fileService.CheckDirId(cf.ID, claimsId)
	if err != nil {
		global.BOX_LOG.Error("文件夹不存在", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = folderService.ChangeFolderName(cf.ID, cf.FolderName)
	if err != nil {
		global.BOX_LOG.Error("修改文件夹出错", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("修改成功", c)
}
