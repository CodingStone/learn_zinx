package ctx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"learn_zinx/Cobra.mayfly/pkg/ginx"
	"learn_zinx/Cobra.mayfly/pkg/model"
	"learn_zinx/Cobra.mayfly/pkg/utils/assert"
	"time"
)

// 处理函数 [这个地方注意 传入函数的处理方法]
type handlerFunc func(*ReqCtx)

type ReqCtx struct {
	GinCtx             *gin.Context        //gin  context
	RequiredPermission *Permission         //需要的权限信息,默认为nil, 需要校验token
	LoginAccount       *model.LoginAccount //登陆账号信息，只有校验token后才会有值

	LogInfo  *LogInfo    //日志相关信息
	ReqParam interface{} //请求参数，主要用于记录日志
	ResData  interface{} //相应结果
	Err      interface{} //请求错误

	timed int64 //执行时间
	NoRes bool  //无需返回结果，即文件下载等
}

// # 执行处理请求核心调用方法
func (rc *ReqCtx) Handle(handler handlerFunc) {
	ginCtx := rc.GinCtx
	defer func() {
		if err := recover(); err != nil {
			rc.Err = err
			ginx.ErrorRes(ginCtx, err)
		}
		//应用所有请求，处理处理器
		ApplyHandlerInterceptor(afterHandlers, rc)
	}()
	// # 上下文一定是被初始化的。
	assert.IsTrue(ginCtx != nil, "ginContext == nil")
	// 默认为不记录请求参数，可在handler回调函数中覆盖赋值
	rc.ReqParam = 0
	//默认相应结果为nil，可在handler中赋值
	rc.ResData = nil

	//调用请求前所有处理器
	err := ApplyHandlerInterceptor(beforeHandlers, rc) //# 调用请求处理前拦截器
	if err != nil {
		panic(err)
	}
	begin := time.Now() // #计算请求耗时

	handler(rc) //# 处理当前请求 相关联的请求上下文
	println("handler 函数退出")
	rc.timed = time.Since(begin).Milliseconds()
	if !rc.NoRes {
		ginx.SuccessRes(ginCtx, rc.ResData) // 返回结果
	}
}

func (rc *ReqCtx) DownLoad(reader io.Reader, filename string) {
	rc.NoRes = true
	ginx.DownLoad(rc.GinCtx, reader, filename)
}

// 新建请求上下文，默认需要校验token
func NewReqCtx() *ReqCtx {
	return &ReqCtx{}
}

// #实例化一个请求请求上下文
func NewReqCtxWithGin(g *gin.Context) *ReqCtx {
	return &ReqCtx{GinCtx: g}
}

// 调用该方法设置请求描述，则默认记录日志，并不记录响应结果
func (r *ReqCtx) WithLog(li *LogInfo) *ReqCtx {
	r.LogInfo = li
	return r
}

// 设置请求上下文需要的权限信息
func (r *ReqCtx) WithRequiredPermission(permission *Permission) *ReqCtx {
	r.RequiredPermission = permission
	return r
}

// 是否需要token [这里表示永远不需要, 注意这块含义]
func (r *ReqCtx) WithNeedToken(needToken bool) *ReqCtx {
	r.RequiredPermission = &Permission{NeedToken: false}
	return r
}

type HandlerInterceptorFunc func(*ReqCtx) error
type HandlerInterceptors []HandlerInterceptorFunc

var (
	beforeHandlers HandlerInterceptors // # 请求前处理器
	afterHandlers  HandlerInterceptors
)

func UseBeforeHandlerInterceptor(b HandlerInterceptorFunc) {
	beforeHandlers = append(beforeHandlers, b)
}

// 使用后置处理器函数
func UseAfterHandlerInterceptor(b HandlerInterceptorFunc) {
	afterHandlers = append(afterHandlers, b)
}

// 应用指定处理器拦截器，如果又一个错误则直接返回错误
func ApplyHandlerInterceptor(his HandlerInterceptors, rc *ReqCtx) interface{} {
	fmt.Printf("HandlerInterceptors len: %v\n", len(his))
	for _, handler := range his {
		if err := handler(rc); err != nil {
			return err
		}
	}
	return nil
}
