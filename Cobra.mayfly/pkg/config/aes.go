package config

import (
	"fmt"
	"learn_zinx/Cobra.mayfly/pkg/utils"
	"learn_zinx/Cobra.mayfly/pkg/utils/assert"
)

// 编码并base64
func (a *Aes) EncryptBase64(data []byte) (string, error) {
	return utils.AesEncryptBase64(data, []byte(a.Key))
}

// base64解码后再aes解码
func (a *Aes) DecryptBase64(data string) ([]byte, error) {
	return utils.AesDecryptBase64(data, []byte(a.Key))
}

type Aes struct {
	Key string `yaml:"key"`
}

func (a *Aes) Valid() {
	aesKeyLen := len(a.Key)
	assert.IsTrue(aesKeyLen == 16 || aesKeyLen == 24 || aesKeyLen == 32, fmt.Sprintf("config.yml之 [aes.key] 长度需为16、24、32位长度, 当前为%d位", aesKeyLen))
}
