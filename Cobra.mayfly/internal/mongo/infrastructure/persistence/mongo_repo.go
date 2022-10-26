package persistence

import (
	"fmt"
	"learn_zinx/Cobra.mayfly/internal/mongo/domain/entity"
	"learn_zinx/Cobra.mayfly/internal/mongo/domain/repository"
	"learn_zinx/Cobra.mayfly/pkg/biz"
	"learn_zinx/Cobra.mayfly/pkg/model"
)

type mongoRepoImpl struct{}

func newMongoRepo() repository.Mongo {
	return new(mongoRepoImpl)
}

// 分页获取数据库信息列表
func (d *mongoRepoImpl) GetList(condition *entity.Mongo, pageParam *model.PageParam, toEntity interface{}, orderBy ...string) *model.PageResult {
	sql := "SELECT d.* FROM t_mongo d JOIN t_project_member pm ON d.project_id = pm.project_id WHERE 1 = 1 "
	if condition.CreatorId != 0 {
		// 使用创建者id模拟项目成员id
		sql = fmt.Sprintf("%s AND pm.account_id = %d", sql, condition.CreatorId)
	}
	if condition.ProjectId != 0 {
		sql = fmt.Sprintf("%s AND d.project_id = %d", sql, condition.ProjectId)
	}
	if condition.EnvId != 0 {
		sql = fmt.Sprintf("%s AND d.env_id = %d", sql, condition.EnvId)
	}
	sql = sql + " ORDER BY d.create_time DESC"
	return model.GetPageBySql(sql, pageParam, toEntity)
}

func (d *mongoRepoImpl) Count(condition *entity.Mongo) int64 {
	return model.CountBy(condition)
}

// 根据条件获取
func (d *mongoRepoImpl) Get(condition *entity.Mongo, cols ...string) error {
	return model.GetBy(condition, cols...)
}

// 根据id获取
func (d *mongoRepoImpl) GetById(id uint64, cols ...string) *entity.Mongo {
	db := new(entity.Mongo)
	if err := model.GetById(db, id, cols...); err != nil {
		return nil

	}
	return db
}

func (d *mongoRepoImpl) Insert(db *entity.Mongo) {
	biz.ErrIsNil(model.Insert(db), "新增mongo信息失败")
}

func (d *mongoRepoImpl) Update(db *entity.Mongo) {
	biz.ErrIsNil(model.UpdateById(db), "更新mongo信息失败")
}

func (d *mongoRepoImpl) Delete(id uint64) {
	model.DeleteById(new(entity.Mongo), id)
}
