package api

import (
	"learn_zinx/Cobra.mayfly/internal/sys/application"
	"learn_zinx/Cobra.mayfly/internal/sys/domain/entity"
	"learn_zinx/Cobra.mayfly/pkg/ctx"
	"learn_zinx/Cobra.mayfly/pkg/ginx"
)

type Syslog struct {
	SyslogApp application.Syslog
}

func (r *Syslog) Syslogs(rc *ctx.ReqCtx) {
	g := rc.GinCtx
	condition := &entity.Syslog{
		Type:      int8(ginx.QueryInt(g, "type", 0)),
		CreatorId: uint64(ginx.QueryInt(g, "creatorId", 0)),
	}
	rc.ResData = r.SyslogApp.GetPageList(condition, ginx.GetPageParam(g), new([]entity.Syslog), "create_time DESC")
}
