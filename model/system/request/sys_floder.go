package request

//创建文件夹
type CreateFolder struct {
	FolderName string `json:"folderName"` //文件夹名称
	ParentId   uint   `json:"folderPath"` //文件夹父ID
}

//修改文件夹
type ChangeFolder struct {
	ID         uint   `json:"id"`
	FolderName string `json:"folderName"` //修改文件名使用
	ParentId   uint   `json:"parentId" `  //移动文件夹可用
}
