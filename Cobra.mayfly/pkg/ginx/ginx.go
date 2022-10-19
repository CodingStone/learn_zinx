package ginx

import (
	"github.com/gin-gonic/gin"
	"io"
	"learn_zinx/Cobra.mayfly/pkg/biz"
	"learn_zinx/Cobra.mayfly/pkg/global"
	"learn_zinx/Cobra.mayfly/pkg/model"
	"net/http"
	"strconv"
)

// 绑定并校验请求结构体参数
func BindJsonAndValid(g *gin.Context, data interface{}) {
	if err := g.ShouldBindJSON(data); err != nil {
		panic(biz.NewBizErr(err.Error()))
	}
}

// 绑定查询字符串到
func BindQuery(g *gin.Context, data interface{}) {
	if err := g.ShouldBindQuery(data); err != nil {
		panic(biz.NewBizErr(err.Error()))
	}
}

// 获取分页参数
func GetPageParam(g *gin.Context) *model.PageParam {
	return &model.PageParam{PageNum: QueryInt(g, "pageNum", 1), PageSize: QueryInt(g, "pageSize", 10)}
}

// 获取查询参数中指定参数值，闭关转化为int
func QueryInt(g *gin.Context, qm string, defaultInt int) int {
	qv := g.Query(qm)
	if qv == "" {
		return defaultInt
	}
	qvi, err := strconv.Atoi(qv)
	biz.ErrIsNil(err, "query param not int")
	return qvi
}

// 获取路径参数
func PathParamInt(g *gin.Context, pm string) int {
	value, _ := strconv.Atoi(g.Param(pm))
	return value
}

// 文件下载
func DownLoad(g *gin.Context, reader io.Reader, filename string) {
	g.Header("Content-Type", "application/octet-stream")
	g.Header("Content-Disposition", "attachment; filename="+filename)
	io.Copy(g.Writer, reader)
}

// 返回统一成功结果
func SuccessRes(g *gin.Context, data interface{}) {
	g.JSON(http.StatusOK, model.Success(data))
}

// 返回失败结果集
func ErrorRes(g *gin.Context, err interface{}) {
	switch t := err.(type) {
	case biz.BizError:
		g.JSON(http.StatusOK, model.Error(t))
	case error:
		g.JSON(http.StatusOK, model.ServerError())
		global.Log.Error(t)
	case string:
		g.JSON(http.StatusOK, model.ServerError())
		global.Log.Error(t)
	default:
		global.Log.Error(t)
	}
}
