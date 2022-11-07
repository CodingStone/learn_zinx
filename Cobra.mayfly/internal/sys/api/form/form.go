package form

// 登录表单 # POST请求的 表单字段
type LoginForm struct {
	Username string `json:"username" binding:"required"`
	Password string `binding:"required"`
	Captcha  string
	Cid      string
}
