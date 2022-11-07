package api

import (
	"learn_zinx/Cobra.mayfly/pkg/captcha"
	"learn_zinx/Cobra.mayfly/pkg/ctx"
)

// # 处理生成验证码操作 【生成的时候不存，验证结果依然交给库取处理】
func GenerateCaptcha(rc *ctx.ReqCtx) {
	id, image := captcha.Generate() //
	rc.ResData = map[string]interface{}{"base64Captcha": image, "cid": id}
}
