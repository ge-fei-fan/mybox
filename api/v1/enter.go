package v1

import "mybox/service"

var (
	userService   = new(service.UserService)
	jwtService    = new(service.JwtService)
	fileService   = new(service.FileUploadAndDownloadService)
	folderService = new(service.FolderService)
)
