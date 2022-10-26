package entity

import "learn_zinx/Cobra.mayfly/pkg/model"

type MachineScript struct {
	model.Model
	Name        string `json:"name"`
	MachineId   uint64 `json:"machineId"` // 机器id
	Type        int    `json:"type"`
	Description string `json:"description"` // 脚本描述
	Params      string `json:"params"`      // 参数列表json
	Script      string `json:"script"`      // 脚本内容
}
