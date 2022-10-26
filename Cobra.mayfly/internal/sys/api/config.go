package api

import (
	"learn_zinx/Cobra.mayfly/internal/sys/api/form"
	"learn_zinx/Cobra.mayfly/internal/sys/application"
	"learn_zinx/Cobra.mayfly/internal/sys/domain/entity"
	"learn_zinx/Cobra.mayfly/pkg/biz"
	"learn_zinx/Cobra.mayfly/pkg/ctx"
	"learn_zinx/Cobra.mayfly/pkg/ginx"
	"learn_zinx/Cobra.mayfly/pkg/utils"
)

type Config struct {
	ConfigApp application.Config
}

func (c *Config) Configs(rc *ctx.ReqCtx) {
	g := rc.GinCtx
	condition := &entity.Config{Key: g.Query("key")}
	rc.ResData = c.ConfigApp.GetPageList(condition, ginx.GetPageParam(g), new([]entity.Config))
}

func (c *Config) GetConfigValueByKey(rc *ctx.ReqCtx) {
	key := rc.GinCtx.Query("key")
	biz.NotEmpty(key, "key不能为空")
	rc.ResData = c.ConfigApp.GetConfig(key).Value
}

func (c *Config) SaveConfig(rc *ctx.ReqCtx) {
	g := rc.GinCtx
	form := &form.ConfigForm{}
	ginx.BindJsonAndValid(g, form)
	rc.ReqParam = form

	config := new(entity.Config)
	utils.Copy(config, form)
	config.SetBaseInfo(rc.LoginAccount)
	c.ConfigApp.Save(config)
}
