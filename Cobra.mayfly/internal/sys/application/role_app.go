package application

import (
	"learn_zinx/Cobra.mayfly/internal/sys/domain/entity"
	"learn_zinx/Cobra.mayfly/internal/sys/domain/repository"
	"learn_zinx/Cobra.mayfly/pkg/model"
	"strings"
)

type Role interface {
	GetPageList(condition *entity.Role, pageParam *model.PageParam, toEntity interface{}, orderBy ...string) *model.PageResult

	SaveRole(role *entity.Role)

	DeleteRole(id uint64)

	GetRoleResourceIds(roleId uint64) []uint64

	GetRoleResources(roleId uint64, toEntity interface{})

	// 保存角色资源关联记录
	SaveRoleResource(rr *entity.RoleResource)

	// 删除角色资源关联记录
	DeleteRoleResource(roleId uint64, resourceId uint64)

	// 获取账号角色id列表
	GetAccountRoleIds(accountId uint64) []uint64

	// 保存账号角色关联信息
	SaveAccountRole(rr *entity.AccountRole)

	DeleteAccountRole(accountId, roleId uint64)

	GetAccountRoles(accountId uint64, toEntity interface{})
}

func newRoleApp(roleRepo repository.Role) Role {
	return &roleAppImpl{
		roleRepo: roleRepo,
	}
}

type roleAppImpl struct {
	roleRepo repository.Role
}

func (m *roleAppImpl) GetPageList(condition *entity.Role, pageParam *model.PageParam, toEntity interface{}, orderBy ...string) *model.PageResult {
	return m.roleRepo.GetPageList(condition, pageParam, toEntity, orderBy...)
}

func (m *roleAppImpl) SaveRole(role *entity.Role) {
	role.Code = strings.ToUpper(role.Code)
	if role.Id != 0 {
		// code不可更改，防止误传
		role.Code = ""
		model.UpdateById(role)
	} else {
		role.Status = 1
		model.Insert(role)
	}
}

func (m *roleAppImpl) DeleteRole(id uint64) {
	m.roleRepo.Delete(id)
	// 删除角色与资源的关联关系
	model.DeleteByCondition(&entity.RoleResource{RoleId: id})
}

func (m *roleAppImpl) GetRoleResourceIds(roleId uint64) []uint64 {
	return m.roleRepo.GetRoleResourceIds(roleId)
}

func (m *roleAppImpl) GetRoleResources(roleId uint64, toEntity interface{}) {
	m.roleRepo.GetRoleResources(roleId, toEntity)
}

func (m *roleAppImpl) SaveRoleResource(rr *entity.RoleResource) {
	m.roleRepo.SaveRoleResource(rr)
}

func (m *roleAppImpl) DeleteRoleResource(roleId uint64, resourceId uint64) {
	m.roleRepo.DeleteRoleResource(roleId, resourceId)
}

func (m *roleAppImpl) GetAccountRoleIds(accountId uint64) []uint64 {
	return m.roleRepo.GetAccountRoleIds(accountId)
}

func (m *roleAppImpl) SaveAccountRole(rr *entity.AccountRole) {
	m.roleRepo.SaveAccountRole(rr)
}

func (m *roleAppImpl) DeleteAccountRole(accountId, roleId uint64) {
	m.roleRepo.DeleteAccountRole(accountId, roleId)
}

func (m *roleAppImpl) GetAccountRoles(accountId uint64, toEntity interface{}) {
	m.roleRepo.GetAccountRoles(accountId, toEntity)
}
