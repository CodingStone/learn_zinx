package api

import (
	"encoding/base64"
	"fmt"
	"learn_zinx/Cobra.mayfly/internal/machine/api/form"
	"learn_zinx/Cobra.mayfly/internal/machine/api/vo"
	"learn_zinx/Cobra.mayfly/internal/machine/application"
	"learn_zinx/Cobra.mayfly/internal/machine/domain/entity"
	"learn_zinx/Cobra.mayfly/internal/machine/infrastructure/machine"
	projectapp "learn_zinx/Cobra.mayfly/internal/project/application"
	"learn_zinx/Cobra.mayfly/pkg/biz"
	"learn_zinx/Cobra.mayfly/pkg/config"
	"learn_zinx/Cobra.mayfly/pkg/ctx"
	"learn_zinx/Cobra.mayfly/pkg/ginx"
	"learn_zinx/Cobra.mayfly/pkg/utils"
	"learn_zinx/Cobra.mayfly/pkg/ws"
	"os"
	"path"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Machine struct {
	MachineApp application.Machine
	ProjectApp projectapp.Project
}

func (m *Machine) Machines(rc *ctx.ReqCtx) {
	condition := new(entity.Machine)
	// 使用创建者id模拟账号成员id
	condition.CreatorId = rc.LoginAccount.Id
	condition.Ip = rc.GinCtx.Query("ip")
	condition.Name = rc.GinCtx.Query("name")
	condition.ProjectId = uint64(ginx.QueryInt(rc.GinCtx, "projectId", 0))

	res := m.MachineApp.GetMachineList(condition, ginx.GetPageParam(rc.GinCtx), new([]*vo.MachineVO))
	if res.Total == 0 {
		rc.ResData = res
		return
	}

	list := res.List.(*[]*vo.MachineVO)
	for _, mv := range *list {
		mv.HasCli = machine.HasCli(*mv.Id)
	}
	rc.ResData = res
}

func (m *Machine) MachineStats(rc *ctx.ReqCtx) {
	stats := m.MachineApp.GetCli(GetMachineId(rc.GinCtx)).GetAllStats()
	rc.ResData = stats
}

func (m *Machine) SaveMachine(rc *ctx.ReqCtx) {
	g := rc.GinCtx
	machineForm := new(form.MachineForm)
	ginx.BindJsonAndValid(g, machineForm)

	me := new(entity.Machine)
	utils.Copy(me, machineForm)

	if me.AuthMethod == entity.MachineAuthMethodPassword {
		// 密码解密，并使用解密后的赋值
		originPwd, err := utils.DefaultRsaDecrypt(machineForm.Password, true)
		biz.ErrIsNilAppendErr(err, "解密密码错误: %s")
		me.Password = originPwd
	}

	// 密码脱敏记录日志
	machineForm.Password = "****"
	rc.ReqParam = machineForm

	me.SetBaseInfo(rc.LoginAccount)
	m.MachineApp.Save(me)
}

// 获取机器实例密码，由于数据库是加密存储，故提供该接口展示原文密码
func (m *Machine) GetMachinePwd(rc *ctx.ReqCtx) {
	mid := GetMachineId(rc.GinCtx)
	me := m.MachineApp.GetById(mid, "Password")
	me.PwdDecrypt()
	rc.ResData = me.Password
}

func (m *Machine) ChangeStatus(rc *ctx.ReqCtx) {
	g := rc.GinCtx
	id := uint64(ginx.PathParamInt(g, "machineId"))
	status := int8(ginx.PathParamInt(g, "status"))
	rc.ReqParam = fmt.Sprintf("id: %d -- status: %d", id, status)
	m.MachineApp.ChangeStatus(id, status)
}

func (m *Machine) DeleteMachine(rc *ctx.ReqCtx) {
	id := uint64(ginx.PathParamInt(rc.GinCtx, "machineId"))
	rc.ReqParam = id
	m.MachineApp.Delete(id)
}

// 关闭机器客户端
func (m *Machine) CloseCli(rc *ctx.ReqCtx) {
	machine.DeleteCli(GetMachineId(rc.GinCtx))
}

// 获取进程列表信息
func (m *Machine) GetProcess(rc *ctx.ReqCtx) {
	g := rc.GinCtx
	cmd := "ps -aux "
	sortType := g.Query("sortType")
	if sortType == "2" {
		cmd += "--sort -pmem "
	} else {
		cmd += "--sort -pcpu "
	}

	pname := g.Query("name")
	if pname != "" {
		cmd += fmt.Sprintf("| grep %s ", pname)
	}

	count := g.Query("count")
	if count == "" {
		count = "10"
	}

	cmd += "| head -n " + count

	cli := m.MachineApp.GetCli(GetMachineId(rc.GinCtx))
	biz.ErrIsNilAppendErr(m.ProjectApp.CanAccess(rc.LoginAccount.Id, cli.GetMachine().ProjectId), "%s")

	res, err := cli.Run(cmd)
	biz.ErrIsNilAppendErr(err, "获取进程信息失败: %s")
	rc.ResData = res
}

// 终止进程
func (m *Machine) KillProcess(rc *ctx.ReqCtx) {
	pid := rc.GinCtx.Query("pid")
	biz.NotEmpty(pid, "进程id不能为空")

	cli := m.MachineApp.GetCli(GetMachineId(rc.GinCtx))
	biz.ErrIsNilAppendErr(m.ProjectApp.CanAccess(rc.LoginAccount.Id, cli.GetMachine().ProjectId), "%s")

	_, err := cli.Run("sudo kill -9 " + pid)
	biz.ErrIsNilAppendErr(err, "终止进程失败: %s")
}

func (m *Machine) WsSSH(g *gin.Context) {
	wsConn, err := ws.Upgrader.Upgrade(g.Writer, g.Request, nil)
	defer func() {
		if wsConn != nil {
			if err := recover(); err != nil {
				wsConn.WriteMessage(websocket.TextMessage, []byte(err.(error).Error()))
			}
			wsConn.Close()
		}
	}()

	biz.ErrIsNilAppendErr(err, "升级websocket失败: %s")
	// 权限校验
	rc := ctx.NewReqCtxWithGin(g).WithRequiredPermission(ctx.NewPermission("machine:terminal"))
	if err = ctx.PermissionHandler(rc); err != nil {
		panic(biz.NewBizErr("\033[1;31m您没有权限操作该机器终端,请重新登录后再试~\033[0m"))
	}

	cli := m.MachineApp.GetCli(GetMachineId(g))
	biz.ErrIsNilAppendErr(m.ProjectApp.CanAccess(rc.LoginAccount.Id, cli.GetMachine().ProjectId), "%s")

	cols := ginx.QueryInt(g, "cols", 80)
	rows := ginx.QueryInt(g, "rows", 40)

	var recorder *machine.Recorder
	if cli.GetMachine().EnableRecorder == 1 {
		now := time.Now()
		// 回放文件路径为: 基础配置路径/机器id/操作日期/操作者账号/操作时间.cast
		recPath := fmt.Sprintf("%s/%d/%s/%s", config.Conf.Server.GetMachineRecPath(), cli.GetMachine().Id, now.Format("20060102"), rc.LoginAccount.Username)
		os.MkdirAll(recPath, 0766)
		fileName := path.Join(recPath, fmt.Sprintf("%s.cast", now.Format("20060102_150405")))
		f, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0766)
		biz.ErrIsNilAppendErr(err, "创建终端回放记录文件失败: %s")
		defer f.Close()
		recorder = machine.NewRecorder(f)
	}

	mts, err := machine.NewTerminalSession(utils.RandString(16), wsConn, cli, rows, cols, recorder)
	biz.ErrIsNilAppendErr(err, "\033[1;31m连接失败: %s\033[0m")
	mts.Start()
	defer mts.Stop()
}

// 获取机器终端回放记录的相应文件夹名或文件内容
func (m *Machine) MachineRecDirNames(rc *ctx.ReqCtx) {
	readPath := rc.GinCtx.Query("path")
	biz.NotEmpty(readPath, "path不能为空")
	path_ := path.Join(config.Conf.Server.GetMachineRecPath(), readPath)

	// 如果是读取文件内容，则读取对应回放记录文件内容，否则读取文件夹名列表。小小偷懒一会不想再加个接口
	isFile := rc.GinCtx.Query("isFile")
	if isFile == "1" {
		bytes, err := os.ReadFile(path_)
		biz.ErrIsNilAppendErr(err, "还未有相应终端操作记录: %s")
		rc.ResData = base64.StdEncoding.EncodeToString(bytes)
		return
	}

	files, err := os.ReadDir(path_)
	biz.ErrIsNilAppendErr(err, "还未有相应终端操作记录: %s")
	var names []string
	for _, f := range files {
		names = append(names, f.Name())
	}
	sort.Sort(sort.Reverse(sort.StringSlice(names)))
	rc.ResData = names
}

func GetMachineId(g *gin.Context) uint64 {
	machineId, _ := strconv.Atoi(g.Param("machineId"))
	biz.IsTrue(machineId != 0, "machineId错误")
	return uint64(machineId)
}
