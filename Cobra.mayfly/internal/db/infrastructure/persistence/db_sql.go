package persistence

import (
	"learn_zinx/Cobra.mayfly/internal/db/domain/entity"
	"learn_zinx/Cobra.mayfly/internal/db/domain/repository"
	"learn_zinx/Cobra.mayfly/pkg/biz"
	"learn_zinx/Cobra.mayfly/pkg/model"
)

type dbSqlRepoImpl struct{}

func newDbSqlRepo() repository.DbSql {
	return new(dbSqlRepoImpl)
}

// 分页获取数据库信息列表
func (d *dbSqlRepoImpl) DeleteBy(condition *entity.DbSql) {
	biz.ErrIsNil(model.DeleteByCondition(condition), "删除sql失败")
}
