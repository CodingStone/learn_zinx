package persistence

import (
	"learn_zinx/Cobra.mayfly/internal/machine/domain/entity"
	"learn_zinx/Cobra.mayfly/internal/machine/domain/repository"
	"learn_zinx/Cobra.mayfly/pkg/biz"
	"learn_zinx/Cobra.mayfly/pkg/model"
)

type machineFileRepoImpl struct{}

func newMachineFileRepo() repository.MachineFile {
	return new(machineFileRepoImpl)
}

// 分页获取机器文件信息列表
func (m *machineFileRepoImpl) GetPageList(condition *entity.MachineFile, pageParam *model.PageParam, toEntity interface{}, orderBy ...string) *model.PageResult {
	return model.GetPage(pageParam, condition, toEntity, orderBy...)
}

// 根据条件获取账号信息
func (m *machineFileRepoImpl) GetMachineFile(condition *entity.MachineFile, cols ...string) error {
	return model.GetBy(condition, cols...)
}

// 根据id获取
func (m *machineFileRepoImpl) GetById(id uint64, cols ...string) *entity.MachineFile {
	ms := new(entity.MachineFile)
	if err := model.GetById(ms, id, cols...); err != nil {
		return nil

	}
	return ms
}

// 根据id获取
func (m *machineFileRepoImpl) Delete(id uint64) {
	biz.ErrIsNil(model.DeleteById(new(entity.MachineFile), id), "删除失败")
}

func (m *machineFileRepoImpl) Create(entity *entity.MachineFile) {
	model.Insert(entity)
}

func (m *machineFileRepoImpl) UpdateById(entity *entity.MachineFile) {
	model.UpdateById(entity)
}
