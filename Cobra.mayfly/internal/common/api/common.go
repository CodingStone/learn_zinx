package api

import (
	"learn_zinx/Cobra.mayfly/pkg/biz"
	"learn_zinx/Cobra.mayfly/pkg/ctx"
	"learn_zinx/Cobra.mayfly/pkg/utils"
)

type Common struct {
}

// # 为客户端生成公钥
func (i *Common) RasPublicKey(rc *ctx.ReqCtx) {
	publicKeyStr, err := utils.GetRsaPublicKey()
	biz.ErrIsNilAppendErr(err, "rsa生成公私钥失败")
	rc.ResData = publicKeyStr
}
