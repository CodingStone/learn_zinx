package persistence

import (
	"fmt"
	"learn_zinx/Cobra.mayfly/internal/redis/domain/entity"
	"learn_zinx/Cobra.mayfly/internal/redis/domain/repository"
	"learn_zinx/Cobra.mayfly/pkg/biz"
	"learn_zinx/Cobra.mayfly/pkg/model"
)

type redisRepoImpl struct{}

func newRedisRepo() repository.Redis {
	return new(redisRepoImpl)
}

// 分页获取机器信息列表
func (r *redisRepoImpl) GetRedisList(condition *entity.Redis, pageParam *model.PageParam, toEntity interface{}, orderBy ...string) *model.PageResult {
	sql := "SELECT d.* FROM t_redis d JOIN t_project_member pm ON d.project_id = pm.project_id WHERE 1 = 1 "
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
	if condition.Host != "" {
		sql = sql + " AND d.host LIKE '%" + condition.Host + "%'"
	}
	sql = sql + " ORDER BY d.create_time DESC"
	return model.GetPageBySql(sql, pageParam, toEntity)
}

func (r *redisRepoImpl) Count(condition *entity.Redis) int64 {
	return model.CountBy(condition)
}

// 根据id获取
func (r *redisRepoImpl) GetById(id uint64, cols ...string) *entity.Redis {
	rd := new(entity.Redis)
	if err := model.GetById(rd, id, cols...); err != nil {
		return nil
	}
	return rd
}

func (r *redisRepoImpl) GetRedis(condition *entity.Redis, cols ...string) error {
	return model.GetBy(condition, cols...)
}

func (r *redisRepoImpl) Insert(redis *entity.Redis) {
	biz.ErrIsNilAppendErr(model.Insert(redis), "新增失败: %s")
}

func (r *redisRepoImpl) Update(redis *entity.Redis) {
	biz.ErrIsNilAppendErr(model.UpdateById(redis), "更新失败: %s")
}

func (r *redisRepoImpl) Delete(id uint64) {
	biz.ErrIsNilAppendErr(model.DeleteById(new(entity.Redis), id), "删除失败: %s")
}
