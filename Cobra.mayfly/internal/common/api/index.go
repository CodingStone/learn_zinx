package api

import (
	"fmt"
	dbapp "learn_zinx/Cobra.mayfly/internal/db/application"
	dbentity "learn_zinx/Cobra.mayfly/internal/db/domain/entity"
	machineapp "learn_zinx/Cobra.mayfly/internal/machine/application"
	machineentity "learn_zinx/Cobra.mayfly/internal/machine/domain/entity"
	projectapp "learn_zinx/Cobra.mayfly/internal/project/application"
	projectentity "learn_zinx/Cobra.mayfly/internal/project/domain/entity"
	redisapp "learn_zinx/Cobra.mayfly/internal/redis/application"
	redisentity "learn_zinx/Cobra.mayfly/internal/redis/domain/entity"
	"learn_zinx/Cobra.mayfly/pkg/ctx"
)

type Index struct {
	ProjectApp projectapp.Project
	MachineApp machineapp.Machine
	DbApp      dbapp.Db
	RedisApp   redisapp.Redis
	Val        int
}

func (i *Index) Count(rc *ctx.ReqCtx) {
	rc.ResData = map[string]interface{}{
		"projectNum": i.ProjectApp.Count(new(projectentity.Project)),
		"machineNum": i.MachineApp.Count(new(machineentity.Machine)),
		"dbNum":      i.DbApp.Count(new(dbentity.Db)),
		"redisNum":   i.RedisApp.Count(new(redisentity.Redis)),
	}
}
func (i *Index) Test(num int) int {
	fmt.Println("I am test!", i.Val+num)
	return i.Val + num
}
