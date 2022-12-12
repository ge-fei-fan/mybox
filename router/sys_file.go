package router

import (
	"github.com/gin-gonic/gin"
	v1 "mybox/api/v1"
)

func InitFileUploadAndDownloadRouter(Router *gin.RouterGroup) {
	fileUploadAndDownloadRouter := Router.Group("file")
	fileUploadAndDownloadApi := new(v1.FileUploadAndDownloadApi)
	{
		fileUploadAndDownloadRouter.POST("upload", fileUploadAndDownloadApi.UploadFile) //上传文件
		fileUploadAndDownloadRouter.DELETE("", fileUploadAndDownloadApi.DeleteFile)     //删除文件
		fileUploadAndDownloadRouter.POST("share", fileUploadAndDownloadApi.ShareFile)   //分享文件
		//fileUploadAndDownloadRouter.GET("shareFileInfo", fileUploadAndDownloadApi.ShareFileInfo) //获取分享信息
		fileUploadAndDownloadRouter.POST("list", fileUploadAndDownloadApi.List)             //文件信息
		fileUploadAndDownloadRouter.POST("renameFile", fileUploadAndDownloadApi.RenameFile) //修改文件名称
		//分片上传
		fileUploadAndDownloadRouter.POST("check", fileUploadAndDownloadApi.Check)
		fileUploadAndDownloadRouter.POST("breakpointcontinue", fileUploadAndDownloadApi.BreakPointContinue)
		fileUploadAndDownloadRouter.POST("breakpointcontinuefinish", fileUploadAndDownloadApi.BreakPointContinueFinish)
	}
	//文件夹
	folderApi := new(v1.FolderApi)
	{
		fileUploadAndDownloadRouter.POST("createFolder", folderApi.CreateFolder) //创建文件夹
		fileUploadAndDownloadRouter.POST("moveFolder", folderApi.MoveFolder)     //创建文件夹
		fileUploadAndDownloadRouter.POST("renameFolder", folderApi.RenameFolder) //修改文件夹名称
	}
}
