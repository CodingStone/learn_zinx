package entity

import "learn_zinx/Cobra.mayfly/pkg/model"

// 项目
type Project struct {
	model.Model
	Name   string `json:"name"`   // 项目名
	Remark string `json:"remark"` // 备注说明
}
