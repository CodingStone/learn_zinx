package persistence

import (
	"learn_zinx/Cobra.mayfly/internal/sys/domain/entity"
	"learn_zinx/Cobra.mayfly/internal/sys/domain/repository"
	"learn_zinx/Cobra.mayfly/pkg/biz"
	"learn_zinx/Cobra.mayfly/pkg/model"
)

type resourceRepoImpl struct{}

func newResourceRepo() repository.Resource {
	return new(resourceRepoImpl)
}

func (r *resourceRepoImpl) GetResourceList(condition *entity.Resource, toEntity interface{}, orderBy ...string) {
	model.ListByOrder(condition, toEntity, orderBy...)
}

func (r *resourceRepoImpl) GetById(id uint64, cols ...string) *entity.Resource {
	res := new(entity.Resource)
	if err := model.GetById(res, id, cols...); err != nil {
		return nil

	}
	return res
}

func (r *resourceRepoImpl) GetByIdIn(ids []uint64, toEntity interface{}, orderBy ...string) {
	model.GetByIdIn(new(entity.Resource), toEntity, ids, orderBy...)
}

func (r *resourceRepoImpl) Delete(id uint64) {
	biz.ErrIsNil(model.DeleteById(new(entity.Resource), id), "删除失败")
}

func (r *resourceRepoImpl) GetByCondition(condition *entity.Resource, cols ...string) error {
	return model.GetBy(condition, cols...)
}

// # 面对复杂查询结构 直接使用sql查询
func (r *resourceRepoImpl) GetAccountResources(accountId uint64, toEntity interface{}) {
	sql := `SELECT
	           m.id,
	           m.pid,
	           m.weight,
	           m.name,
	           m.code,
	           m.meta,
	           m.type,
	           m.status
            FROM
	           t_sys_resource m 
            WHERE
         	   m.STATUS = 1 
	        AND m.id IN (
	            SELECT DISTINCT
		            ( rmb.resource_id ) 
	            FROM
		            t_sys_account_role p
		        JOIN t_sys_role r ON p.role_Id = r.id 
		        AND p.account_id = ?
		        AND r.STATUS = 1
		        JOIN t_sys_role_resource rmb ON rmb.role_id = r.id UNION
	            SELECT
		            r.id 
	            FROM
	             	t_sys_resource r
		        JOIN t_sys_role_resource rr ON r.id = rr.resource_id
		        JOIN t_sys_role ro ON rr.role_id = ro.id 
		        AND ro.status = 1 AND ro.code LIKE 'COMMON%' 
	        ) 
            ORDER BY
	        m.pid ASC,
	        m.weight ASC`
	biz.ErrIsNilAppendErr(model.GetListBySql2Model(sql, toEntity, accountId), "查询账号资源失败: %s")
}
