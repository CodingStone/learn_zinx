package persistence

import (
	"fmt"
	"learn_zinx/Cobra.mayfly/internal/machine/domain/entity"
	"learn_zinx/Cobra.mayfly/internal/machine/domain/repository"
	"learn_zinx/Cobra.mayfly/pkg/biz"
	"learn_zinx/Cobra.mayfly/pkg/model"
)

type machineRepoImpl struct{}

func newMachineRepo() repository.Machine {
	return new(machineRepoImpl)
}

// 分页获取机器信息列表
func (m *machineRepoImpl) GetMachineList(condition *entity.Machine, pageParam *model.PageParam, toEntity interface{}, orderBy ...string) *model.PageResult {
	sql := "SELECT m.* FROM t_machine m JOIN t_project_member pm ON m.project_id = pm.project_id WHERE 1 = 1 "
	if condition.CreatorId != 0 {
		// 使用创建者id模拟项目成员id
		sql = fmt.Sprintf("%s AND pm.account_id = %d", sql, condition.CreatorId)
	}
	if condition.ProjectId != 0 {
		sql = fmt.Sprintf("%s AND m.project_id = %d", sql, condition.ProjectId)
	}
	if condition.Ip != "" {
		sql = sql + " AND m.ip LIKE '%" + condition.Ip + "%'"
	}
	if condition.Name != "" {
		sql = sql + " AND m.name LIKE '%" + condition.Name + "%'"
	}
	sql = sql + " ORDER BY m.project_id, m.create_time DESC"
	return model.GetPageBySql(sql, pageParam, toEntity)
}

func (m *machineRepoImpl) Count(condition *entity.Machine) int64 {
	return model.CountBy(condition)
}

// 根据条件获取账号信息
func (m *machineRepoImpl) GetMachine(condition *entity.Machine, cols ...string) error {
	return model.GetBy(condition, cols...)
}

// 根据id获取
func (m *machineRepoImpl) GetById(id uint64, cols ...string) *entity.Machine {
	machine := new(entity.Machine)
	if err := model.GetById(machine, id, cols...); err != nil {
		return nil

	}
	return machine
}

func (m *machineRepoImpl) Create(entity *entity.Machine) {
	biz.ErrIsNilAppendErr(model.Insert(entity), "创建机器信息失败: %s")
}

func (m *machineRepoImpl) UpdateById(entity *entity.Machine) {
	biz.ErrIsNilAppendErr(model.UpdateById(entity), "更新机器信息失败: %s")
}
