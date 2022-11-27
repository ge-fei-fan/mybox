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
		fileUploadAndDownloadRouter.POST("list", fileUploadAndDownloadApi.List)         //文件信息
	}
	//文件夹
	folderApi := new(v1.FolderApi)
	{
		fileUploadAndDownloadRouter.POST("createFolder", folderApi.CreateFolder)         //创建文件夹
		fileUploadAndDownloadRouter.POST("moveFolder", folderApi.MoveFolder)             //创建文件夹
		fileUploadAndDownloadRouter.POST("changeFolderName", folderApi.ChangeFolderName) //修改文件夹名称
	}
}
