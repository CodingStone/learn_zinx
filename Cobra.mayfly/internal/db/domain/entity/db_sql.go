package entity

import (
	"learn_zinx/Cobra.mayfly/pkg/model"
)

type DbSql struct {
	model.Model `orm:"-"`

	DbId uint64 `json:"dbId"`
	Db   string `json:"db"`
	Type int    `json:"type"` // 类型
	Sql  string `json:"sql"`
	Name string `json:"name"`
}
