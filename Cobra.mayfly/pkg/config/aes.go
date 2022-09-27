package config

import (
	"fmt"
	"learn_zinx/Cobra.mayfly/pkg/utils/assert"
)

type Aes struct {
	Key string `yaml:"key"`
}

func (a *Aes) Valid() {
	aesKeyLen := len(a.Key)
	assert.IsTrue(aesKeyLen == 16 || aesKeyLen == 24 || aesKeyLen == 32, fmt.Sprintf("config.yml之 [aes.key] 长度需为16、24、32位长度, 当前为%d位", aesKeyLen))
}
