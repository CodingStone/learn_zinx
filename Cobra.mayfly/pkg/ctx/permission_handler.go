package ctx

import (
	"fmt"
	"learn_zinx/Cobra.mayfly/pkg/biz"
	"learn_zinx/Cobra.mayfly/pkg/cache"
	"learn_zinx/Cobra.mayfly/pkg/config"
	"time"
)

type Permission struct {
	NeedToken bool   //是否需要token
	Code      string //权限Code
}

func NewPermission(code string) *Permission {
	return &Permission{NeedToken: true, Code: code}
}

func (p *Permission) WithNeedToken(needToken bool) *Permission {
	p.NeedToken = needToken
	return p
}

type PermissionCodeRegistry interface {
	// 保存用户权限code
	SaveCodes(userId uint64, codes []string)

	// 判断用户是否拥有该code的权限
	HasCode(userId uint64, code string) bool

	Remove(userId uint64)
}

type DefaultPermissionCodeRegistry struct {
	cache *cache.TimedCache
}

// # 保存权限缓存
func (r *DefaultPermissionCodeRegistry) SaveCodes(userId uint64, codes []string) {
	if r.cache == nil {
		// #cache是自己实现的
		r.cache = cache.NewTimedCache(time.Minute*time.Duration(config.Conf.Jwt.ExpireTime), 5*time.Second)
	}
	r.cache.Put(fmt.Sprintf("%v", userId), codes)
}

func (r *DefaultPermissionCodeRegistry) HasCode(userId uint64, code string) bool {
	if r.cache == nil {
		return false
	}
	codes, found := r.cache.Get(fmt.Sprintf("%v", userId))
	if !found {
		return false
	}
	for _, v := range codes.([]string) {
		if v == code {
			return true
		}
	}
	return false
}

func (r *DefaultPermissionCodeRegistry) Remove(userId uint64) {
	r.cache.Delete(fmt.Sprintf("%v", userId))
}

var (
	permissionCodeRegistry PermissionCodeRegistry = &DefaultPermissionCodeRegistry{}
	// permissionError                               = biz.NewBizErrCode(biz.TokenErrorCode, biz.TokenErrorMsg)
)

// 保存用户权限code
func SavePermissionCodes(userId uint64, codes []string) {
	permissionCodeRegistry.SaveCodes(userId, codes)
}

// 删除用户权限code
func DeletePermissionCodes(userId uint64) {
	permissionCodeRegistry.Remove(userId)
}

// 设置权限code注册器
func SetPermissionCodeRegistery(pcr PermissionCodeRegistry) {
	permissionCodeRegistry = pcr
}

func PermissionHandler(rc *ReqCtx) error {
	permission := rc.RequiredPermission
	// 如果需要的权限信息不为空，并且不需要token，则不返回错误，继续后续逻辑
	// fmt.Printf("权限验证信息: %v, code:%s\n", permission.NeedToken, permission.Code)
	if permission != nil && !permission.NeedToken {
		return nil
	}
	tokenStr := rc.GinCtx.Request.Header.Get("Authorization")
	//header 不存在则从token中查询
	if tokenStr == "" {
		tokenStr = rc.GinCtx.Query("token")
	}
	if tokenStr == "" {
		return biz.PermissionErr
	}
	loginAccount, err := ParseToken(tokenStr)
	if err != nil || loginAccount == nil {
		return biz.PermissionErr
	}
	// 权限不为nil，并且permission code不为空，则校验是否有权限code
	if permission != nil && permission.Code != "" {
		if !permissionCodeRegistry.HasCode(loginAccount.Id, permission.Code) {
			return biz.PermissionErr
		}
	}
	// # 在token验证里面 获得 用户信息
	rc.LoginAccount = loginAccount
	return nil
}
