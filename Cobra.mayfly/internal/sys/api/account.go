package api

import (
	"fmt"
	"learn_zinx/Cobra.mayfly/internal/sys/api/form"
	"learn_zinx/Cobra.mayfly/internal/sys/api/vo"
	"learn_zinx/Cobra.mayfly/internal/sys/application"
	"learn_zinx/Cobra.mayfly/internal/sys/domain/entity"
	"learn_zinx/Cobra.mayfly/pkg/biz"
	"learn_zinx/Cobra.mayfly/pkg/ctx"
	"learn_zinx/Cobra.mayfly/pkg/ginx"
	"learn_zinx/Cobra.mayfly/pkg/model"
	"learn_zinx/Cobra.mayfly/pkg/utils"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Account struct {
	AccountApp  application.Account
	ResourceApp application.Resource
	RoleApp     application.Role
	MsgApp      application.Msg
	ConfigApp   application.Config
}

/**   登录者个人相关操作   **/

// @router /accounts/login [post]
func (a *Account) Login(rc *ctx.ReqCtx) {
	loginForm := &form.LoginForm{}              // # 获得表单数据，并将数据赋值给特定值的
	ginx.BindJsonAndValid(rc.GinCtx, loginForm) // # 验证值类型

	// 判断是否有开启登录验证码校验
	//if a.ConfigApp.GetConfig(entity.ConfigKeyUseLoginCaptcha).BoolValue(true) { // # 从db中判断是不是需要验证码
	//	// 校验验证码
	//	biz.IsTrue(captcha.Verify(loginForm.Cid, loginForm.Captcha), "验证码错误") // # 用的Cid（密钥生成id 和 验证码去验证）
	//}

	// # 用于解密获得原始密码
	originPwd, err := utils.DefaultRsaDecrypt(loginForm.Password, true)
	biz.ErrIsNilAppendErr(err, "解密密码错误: %s")

	account := &entity.Account{Username: loginForm.Username}
	fmt.Printf("account:%+v\n", account)
	err = a.AccountApp.GetAccount(account, "Id", "Username", "Password", "Status", "LastLoginTime", "LastLoginIp")
	biz.ErrIsNil(err, "用户名或密码错误(查询错误)")
	fmt.Printf("originPwd is: %v, %v\n", originPwd, account.Password)
	biz.IsTrue(utils.CheckPwdHash(originPwd, account.Password), "用户名或密码错误")
	biz.IsTrue(account.IsEnable(), "该账号不可用")

	// 校验密码强度是否符合
	biz.IsTrueBy(CheckPasswordLever(originPwd), biz.NewBizErrCode(401, "您的密码安全等级较低，请修改后重新登录"))

	var resources vo.AccountResourceVOList
	// 获取账号菜单资源
	a.ResourceApp.GetAccountResources(account.Id, &resources)
	// 菜单树与权限code数组
	var menus vo.AccountResourceVOList
	var permissions []string
	for _, v := range resources {
		if v.Type == entity.ResourceTypeMenu {
			menus = append(menus, v)
		} else {
			permissions = append(permissions, *v.Code)
		}
	}
	// 保存该账号的权限codes
	ctx.SavePermissionCodes(account.Id, permissions)

	clientIp := rc.GinCtx.ClientIP()
	// 保存登录消息
	go a.saveLogin(account, clientIp)

	rc.ReqParam = fmt.Sprintln("登录ip: ", clientIp)
	// 赋值loginAccount 主要用于记录操作日志，因为操作日志保存请求上下文没有该信息不保存日志
	rc.LoginAccount = &model.LoginAccount{Id: account.Id, Username: account.Username}

	rc.ResData = map[string]interface{}{
		"token":         ctx.CreateToken(account.Id, account.Username),
		"username":      account.Username,
		"lastLoginTime": account.LastLoginTime,
		"lastLoginIp":   account.LastLoginIp,
		"menus":         menus.ToTrees(0),
		"permissions":   permissions,
	}
}

func (a *Account) ChangePassword(rc *ctx.ReqCtx) {
	form := new(form.AccountChangePasswordForm)
	ginx.BindJsonAndValid(rc.GinCtx, form)

	originOldPwd, err := utils.DefaultRsaDecrypt(form.OldPassword, true)
	biz.ErrIsNilAppendErr(err, "解密旧密码错误: %s")

	account := &entity.Account{Username: form.Username}
	err = a.AccountApp.GetAccount(account, "Id", "Username", "Password", "Status")
	biz.ErrIsNil(err, "旧密码错误")
	biz.IsTrue(utils.CheckPwdHash(originOldPwd, account.Password), "旧密码错误")
	biz.IsTrue(account.IsEnable(), "该账号不可用")

	originNewPwd, err := utils.DefaultRsaDecrypt(form.NewPassword, true)
	biz.ErrIsNilAppendErr(err, "解密新密码错误: %s")
	biz.IsTrue(CheckPasswordLever(originNewPwd), "密码强度必须8位以上且包含字⺟⼤⼩写+数字+特殊符号")

	updateAccount := new(entity.Account)
	updateAccount.Id = account.Id
	updateAccount.Password = utils.PwdHash(originNewPwd)
	a.AccountApp.Update(updateAccount)

	// 赋值loginAccount 主要用于记录操作日志，因为操作日志保存请求上下文没有该信息不保存日志
	rc.LoginAccount = &model.LoginAccount{Id: account.Id, Username: account.Username}
}

func CheckPasswordLever(ps string) bool {
	if len(ps) < 8 {
		return false
	}
	num := `[0-9]{1}`
	a_z := `[a-zA-Z]{1}`
	symbol := `[!@#~$%^&*()+|_.,]{1}`
	if b, err := regexp.MatchString(num, ps); !b || err != nil {
		return false
	}
	if b, err := regexp.MatchString(a_z, ps); !b || err != nil {
		return false
	}
	if b, err := regexp.MatchString(symbol, ps); !b || err != nil {
		return false
	}
	return true
}

// 保存更新账号登录信息
func (a *Account) saveLogin(account *entity.Account, ip string) {
	// 更新账号最后登录时间
	now := time.Now()
	updateAccount := &entity.Account{LastLoginTime: &now}
	updateAccount.Id = account.Id
	updateAccount.LastLoginIp = ip
	a.AccountApp.Update(updateAccount)

	// 创建登录消息
	loginMsg := &entity.Msg{
		RecipientId: int64(account.Id),
		Msg:         fmt.Sprintf("于%s登录", now.Format("2006-01-02 15:04:05")),
		Type:        1,
	}
	loginMsg.CreateTime = &now
	loginMsg.Creator = account.Username
	loginMsg.CreatorId = account.Id
	a.MsgApp.Create(loginMsg)
}

// 获取个人账号信息
func (a *Account) AccountInfo(rc *ctx.ReqCtx) {
	ap := new(vo.AccountPersonVO)
	// 角色信息
	roles := new([]vo.AccountRoleVO)
	a.RoleApp.GetAccountRoles(rc.LoginAccount.Id, roles)

	ap.Roles = *roles
	rc.ResData = ap
}

// 更新个人账号信息
func (a *Account) UpdateAccount(rc *ctx.ReqCtx) {
	updateForm := &form.AccountUpdateForm{}
	ginx.BindJsonAndValid(rc.GinCtx, updateForm)

	updateAccount := new(entity.Account)
	utils.Copy(updateAccount, updateForm)
	// 账号id为登录者账号
	updateAccount.Id = rc.LoginAccount.Id

	if updateAccount.Password != "" {
		biz.IsTrue(CheckPasswordLever(updateAccount.Password), "密码强度必须8位以上且包含字⺟⼤⼩写+数字+特殊符号")
		updateAccount.Password = utils.PwdHash(updateAccount.Password)
	}
	a.AccountApp.Update(updateAccount)
}

// 获取账号接收的消息列表
func (a *Account) GetMsgs(rc *ctx.ReqCtx) {
	condition := &entity.Msg{
		RecipientId: int64(rc.LoginAccount.Id),
	}
	rc.ResData = a.MsgApp.GetPageList(condition, ginx.GetPageParam(rc.GinCtx), new([]entity.Msg))
}

/**    后台账号操作    **/

// @router /accounts [get]
func (a *Account) Accounts(rc *ctx.ReqCtx) {
	condition := &entity.Account{}
	condition.Username = rc.GinCtx.Query("username")
	rc.ResData = a.AccountApp.GetPageList(condition, ginx.GetPageParam(rc.GinCtx), new([]vo.AccountManageVO))
}

// @router /accounts
func (a *Account) SaveAccount(rc *ctx.ReqCtx) {
	form := &form.AccountCreateForm{}
	ginx.BindJsonAndValid(rc.GinCtx, form)
	rc.ReqParam = form

	account := &entity.Account{}
	utils.Copy(account, form)
	account.SetBaseInfo(rc.LoginAccount)

	if account.Id == 0 {
		a.AccountApp.Create(account)
	} else {
		if account.Password != "" {
			biz.IsTrue(CheckPasswordLever(account.Password), "密码强度必须8位以上且包含字⺟⼤⼩写+数字+特殊符号")
			account.Password = utils.PwdHash(account.Password)
		}
		a.AccountApp.Update(account)
	}
}

func (a *Account) ChangeStatus(rc *ctx.ReqCtx) {
	g := rc.GinCtx

	account := &entity.Account{}
	account.Id = uint64(ginx.PathParamInt(g, "id"))
	account.Status = int8(ginx.PathParamInt(g, "status"))
	rc.ReqParam = fmt.Sprintf("accountId: %d, status: %d", account.Id, account.Status)
	a.AccountApp.Update(account)
}

func (a *Account) DeleteAccount(rc *ctx.ReqCtx) {
	id := uint64(ginx.PathParamInt(rc.GinCtx, "id"))
	rc.ReqParam = id
	a.AccountApp.Delete(id)
}

// 获取账号角色id列表，用户回显角色分配
func (a *Account) AccountRoleIds(rc *ctx.ReqCtx) {
	rc.ResData = a.RoleApp.GetAccountRoleIds(uint64(ginx.PathParamInt(rc.GinCtx, "id")))
}

// 获取账号角色id列表，用户回显角色分配
func (a *Account) AccountRoles(rc *ctx.ReqCtx) {
	vos := new([]vo.AccountRoleVO)
	a.RoleApp.GetAccountRoles(uint64(ginx.PathParamInt(rc.GinCtx, "id")), vos)
	rc.ResData = vos
}

func (a *Account) AccountResources(rc *ctx.ReqCtx) {
	var resources vo.ResourceManageVOList
	// 获取账号菜单资源
	a.ResourceApp.GetAccountResources(uint64(ginx.PathParamInt(rc.GinCtx, "id")), &resources)
	rc.ResData = resources.ToTrees(0)
}

// 保存账号角色信息
func (a *Account) SaveRoles(rc *ctx.ReqCtx) {
	g := rc.GinCtx

	var form form.AccountRoleForm
	ginx.BindJsonAndValid(g, &form)
	aid := uint64(form.Id)
	rc.ReqParam = form

	// 将,拼接的字符串进行切割
	idsStr := strings.Split(form.RoleIds, ",")
	var newIds []interface{}
	for _, v := range idsStr {
		id, _ := strconv.Atoi(v)
		newIds = append(newIds, uint64(id))
	}

	// 将[]uint64转为[]interface{}
	oIds := a.RoleApp.GetAccountRoleIds(uint64(form.Id))
	var oldIds []interface{}
	for _, v := range oIds {
		oldIds = append(oldIds, v)
	}

	addIds, delIds, _ := utils.ArrayCompare(newIds, oldIds, func(i1, i2 interface{}) bool {
		return i1.(uint64) == i2.(uint64)
	})

	createTime := time.Now()
	creator := rc.LoginAccount.Username
	creatorId := rc.LoginAccount.Id
	for _, v := range addIds {
		rr := &entity.AccountRole{AccountId: aid, RoleId: v.(uint64), CreateTime: &createTime, CreatorId: creatorId, Creator: creator}
		a.RoleApp.SaveAccountRole(rr)
	}
	for _, v := range delIds {
		a.RoleApp.DeleteAccountRole(aid, v.(uint64))
	}
}
