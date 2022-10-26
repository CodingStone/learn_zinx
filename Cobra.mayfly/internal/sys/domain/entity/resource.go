package entity

import "learn_zinx/Cobra.mayfly/pkg/model"

type Resource struct {
	model.Model
	Pid    int    `json:"pid"`
	Type   int8   `json:"type"`   // 1：菜单路由；2：资源（按钮等）
	Status int8   `json:"status"` // 1：可用；-1：不可用
	Code   string `json:"code"`
	Name   string `json:"name"`
	Weight int    `json:"weight"`
	Meta   string `json:"meta"`
}

func (a *Resource) TableName() string {
	return "t_sys_resource"
}

const (
	ResourceStatusEnable  int8 = 1  // 启用状态
	ResourceStatusDisable int8 = -1 // 禁用状态

	// 资源状态
	ResourceTypeMenu       int8 = 1
	ResourceTypePermission int8 = 2
)
