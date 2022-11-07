package api

import (
	"fmt"
	"learn_zinx/Cobra.mayfly/internal/sys/api/form"
	"learn_zinx/Cobra.mayfly/internal/sys/application"
	"learn_zinx/Cobra.mayfly/internal/sys/domain/entity"
	"learn_zinx/Cobra.mayfly/pkg/biz"
	"learn_zinx/Cobra.mayfly/pkg/ctx"
	"learn_zinx/Cobra.mayfly/pkg/ginx"
	"learn_zinx/Cobra.mayfly/pkg/utils"
)

// # 与具体业务关联起来 ConfigApp 是数据库链接。 并实现基本操作
type Config struct {
	ConfigApp application.Config
}

func (c *Config) Configs(rc *ctx.ReqCtx) {
	g := rc.GinCtx
	condition := &entity.Config{Key: g.Query("key")}
	rc.ResData = c.ConfigApp.GetPageList(condition, ginx.GetPageParam(g), new([]entity.Config))
}

func (c *Config) GetConfigValueByKey(rc *ctx.ReqCtx) {
	key := rc.GinCtx.Query("key") // # 获得请求中key参数
	fmt.Printf("key is: %s\n", key)
	biz.NotEmpty(key, "key不能为空")
	// #具体
	rc.ResData = c.ConfigApp.GetConfig(key).Value // # 进行数据库操作 查询得到数据
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
