package response

import "mybox/model/system"

type FileResponse struct {
	File system.FileUploadAndDownload `json:"file"`
}

type ListItem struct {
	Folder []system.SysUserRepository     `json:"folder" ` //目录
	File   []system.FileUploadAndDownload `json:"file"`    //文件
}
