package persistence

import (
	"learn_zinx/Cobra.mayfly/internal/sys/domain/entity"
	"learn_zinx/Cobra.mayfly/internal/sys/domain/repository"
	"learn_zinx/Cobra.mayfly/pkg/biz"
	"learn_zinx/Cobra.mayfly/pkg/model"
)

type accountRepoImpl struct{}

func newAccountRepo() repository.Account {
	return new(accountRepoImpl)
}

func (a *accountRepoImpl) GetAccount(condition *entity.Account, cols ...string) error {
	return model.GetBy(condition, cols...)
}

func (m *accountRepoImpl) GetPageList(condition *entity.Account, pageParam *model.PageParam, toEntity interface{}, orderBy ...string) *model.PageResult {
	sql := "SELECT * FROM t_sys_account "
	username := condition.Username
	if username != "" {
		sql = sql + " WHERE username LIKE '%" + username + "%'"
	}
	return model.GetPageBySql(sql, pageParam, toEntity)
	// return model.GetPage(pageParam, condition, toEntity, orderBy...)
}

func (m *accountRepoImpl) Insert(account *entity.Account) {
	biz.ErrIsNil(model.Insert(account), "新增账号信息失败")
}

func (m *accountRepoImpl) Update(account *entity.Account) {
	biz.ErrIsNil(model.UpdateById(account), "更新账号信息失败")
}
