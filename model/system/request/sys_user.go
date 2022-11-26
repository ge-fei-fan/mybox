package request

// Register User register structure
type Register struct {
	Username  string `json:"userName" example:"用户名"`
	Password  string `json:"passWord" example:"密码"`
	HeaderImg string `json:"headerImg" example:"头像链接"`
	Email     string `json:"email" example:"电子邮箱"`
}

// User login structure
type Login struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
	//Captcha   string `json:"captcha"`   // 验证码
	//CaptchaId string `json:"captchaId"` // 验证码ID
}
