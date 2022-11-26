package upload

import (
	"mime/multipart"
	"mybox/global"
)

type OSS interface {
	//返回值 访问路径，文件名，错误信息
	UploadFile(file *multipart.FileHeader) (string, string, error)
	DeleteFile(key string) error
}

func NewOss() OSS {
	switch global.BOX_CONFIG.System.OssType {
	case "local":
		return &Local{}
	//case "qiniu":
	//	return &Qiniu{}
	//case "tencent-cos":
	//	return &TencentCOS{}
	//case "aliyun-oss":
	//	return &AliyunOSS{}
	//case "huawei-obs":
	//	return HuaWeiObs
	//case "aws-s3":
	//	return &AwsS3{}
	default:
		return &Local{}
	}
}
