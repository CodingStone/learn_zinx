package ctx

import "learn_zinx/Cobra.mayfly/pkg/config"

var (
	JwtKey  string
	ExpTime uint64
)

func InitTokenConfig() {
	JwtKey = config.Conf.Jwt.Key
	ExpTime = config.Conf.Jwt.ExpireTime
}
