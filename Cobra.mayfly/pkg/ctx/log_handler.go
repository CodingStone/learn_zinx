package ctx

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"learn_zinx/Cobra.mayfly/pkg/biz"
	"learn_zinx/Cobra.mayfly/pkg/logger"
	"learn_zinx/Cobra.mayfly/pkg/utils"
	"reflect"
	"runtime/debug"
)

type SaveLogFunc func(ctx *ReqCtx)

var saveLog SaveLogFunc

// 设置保存日志处理函数
func SetSaveLogFunc(sl SaveLogFunc) {
	saveLog = sl
}

type LogInfo struct {
	LogResp     bool   //是否记录返回结果
	Description string //请求描述
	Save        bool   //是否保存日志
}

// 新建日志信息
func NewLogInfo(description string) *LogInfo {
	return &LogInfo{Description: description, LogResp: false}
}

// 是否记录返回结果
func (i *LogInfo) WithLogResp(logResp bool) *LogInfo {
	i.LogResp = logResp
	return i
}

// 是否保存日志
func (i *LogInfo) WithSave(saveLog bool) *LogInfo {
	i.Save = saveLog
	return i
}

func LogHandler(rc *ReqCtx) error {
	li := rc.LogInfo
	if li == nil {
		return nil
	}
	lfs := logrus.Fields{}
	// 保存登陆账号信息
	if la := rc.LoginAccount; la != nil {
		lfs["uid"] = la.Id
		lfs["uname"] = la.Username
	}
	req := rc.GinCtx.Request
	lfs[req.Method] = req.URL.Path

	//如果需要保存日志，并且保存日志处理函数存在则执行保存日志函数
	if li.Save && saveLog != nil {
		go saveLog(rc)
	}
	if err := rc.Err; err != nil {
		logger.Log.WithFields(lfs).Error(getErrMsg(rc, err))
		return nil
	}
	logger.Log.WithFields(lfs).Info(getLogMsg(rc))
	return nil
}

func getLogMsg(rc *ReqCtx) string {
	msg := rc.LogInfo.Description + fmt.Sprintf(" ->%dms", rc.timed)
	if !utils.IsBlank(reflect.ValueOf(rc.ReqParam)) {
		rb, _ := json.Marshal(rc.ReqParam)
		msg = fmt.Sprintf("\n--> %s", string(rb))
	}
	if rc.LogInfo.LogResp && !utils.IsBlank(reflect.ValueOf(rc.ResData)) {
		respB, _ := json.Marshal(rc.ResData)
		msg = msg + fmt.Sprintf("\n<-- %s", string(respB))
	}
	return msg
}

func getErrMsg(rc *ReqCtx, err interface{}) string {
	msg := rc.LogInfo.Description
	if !utils.IsBlank(reflect.ValueOf(rc.ReqParam)) {
		rb, _ := json.Marshal(rc.ReqParam)
		msg = msg + fmt.Sprintf("\n--> %s", string(rb))
	}
	var errMsg string
	switch t := err.(type) {
	case biz.BizError:
		errMsg = fmt.Sprintf("\n<-e errCode: %d, errMsg: %s", t.Code(), t.Error())
	case error:
		errMsg = fmt.Sprintf("\n<-e errMsg: %s\n%s", t.Error(), string(debug.Stack()))
	case string:
		errMsg = fmt.Sprintf("\n<-e errMsg: %s\n%s", t, string(debug.Stack()))
	}
	return (msg + errMsg)
}
