package ctx

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"learn_zinx/Cobra.mayfly/pkg/biz"
	"learn_zinx/Cobra.mayfly/pkg/config"
	"learn_zinx/Cobra.mayfly/pkg/global"
	"learn_zinx/Cobra.mayfly/pkg/model"
	"learn_zinx/Cobra.mayfly/pkg/utils"
	"time"
)

var (
	JwtKey  string
	ExpTime uint64
)

func InitTokenConfig() {
	JwtKey = config.Conf.Jwt.Key
	ExpTime = config.Conf.Jwt.ExpireTime
}

func CreateToken(userId uint64, username string) string {
	// 带权限创建令牌
	// 设置有效期，过期需要重新登录获取token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       userId,
		"username": username,
		"exp":      time.Now().Add(time.Minute * time.Duration(ExpTime)).Unix(),
	})
	//如果配置文件中的jwt key为空，则随机生成字符串
	if JwtKey == "" {
		JwtKey = utils.RandString(32)
		global.Log.Infof("config.yml未配置jwt.key, 随机生成key为: %s", JwtKey)
	}
	// 使用自定义字符串加密 and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(JwtKey))
	biz.ErrIsNil(err, "token创建失败")
	return tokenString
}

func ParseToken(tokenStr string) (*model.LoginAccount, error) {
	if tokenStr == "" {
		return nil, errors.New("token error")
	}
	// Parse token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtKey), nil
	})
	if err != nil || token == nil {
		return nil, err
	}
	i := token.Claims.(jwt.MapClaims)
	return &model.LoginAccount{Id: uint64(i["id"].(float64)), Username: i["username"].(string)}, nil
}
