package initialize

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/fs"
	"learn_zinx/Cobra.mayfly/pkg/config"
	"learn_zinx/Cobra.mayfly/pkg/middleware"
	"learn_zinx/Cobra.mayfly/static"
	"net/url"
	"reflect"

	"net/http"
)

type Message struct {
	Channel string `uri:"channel" binding:"required"`
	Action  string `uri:"action" binding:"required"`
}
type Data struct {
	Message
	Query url.Values
}

func WrapStaticHandler(h http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Cahe-control", `max-age=31536000`)
		h.ServeHTTP(c.Writer, c.Request)
	}
}
func InitRouter() *gin.Engine {
	//server 配置
	serverConfig := config.Conf.Server
	gin.SetMode(serverConfig.Model)

	var router = gin.New()
	router.MaxMultipartMemory = 8 << 20 // 用于限制上传文件的大小  8M，Gin框架默认是32M

	// 处理路由不存在情况
	router.NoRoute(func(g *gin.Context) {
		g.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": fmt.Sprintf("not found '%s:%s'", g.Request.Method, g.Request.URL.Path)})
	})

	// 使用embed打包静态资源至二进制文件中
	fsys, _ := fs.Sub(static.Static, "static")
	fileServer := http.FileServer(http.FS(fsys))
	handler := WrapStaticHandler(fileServer)
	router.GET("/", handler)
	router.GET("/favicon.ico", handler)
	router.GET("/config.js", handler)
	// 所有/assets/**开头的都是静态资源文件
	router.GET("/assets/*file", handler)

	// 设置静态资源
	if staticConfs := serverConfig.Static; staticConfs != nil {
		for _, scs := range *staticConfs {
			router.StaticFS(scs.RelativePath, http.Dir(scs.Root))
		}
	}

	// 设置静态文件
	if staticFileConfs := serverConfig.StaticFile; staticFileConfs != nil {
		for _, sfs := range *staticFileConfs {
			router.StaticFile(sfs.RelativePath, sfs.Filepath)
		}
	}

	// 是否允许跨域
	if serverConfig.Cors {
		router.Use(middleware.Cors())
	}
	// 设置路由组,  原来方法
	//api := router.Group("/api")
	//{
	//	common_router.InitIndexRouter(api)
	//	common_router.InitCommonRouter(api)
	//	sys_router.Init(api) // # 这里会注册系统路由，包括登陆等多个路由
	//	project_router.Init(api)
	//	machine_router.Init(api)
	//	// 这里面还有很多router
	//
	//}
	router.GET("/api/:channel/:action", func(c *gin.Context) {
		var message Message

		if err := c.ShouldBindUri(&message); err != nil {
			c.JSON(400, gin.H{"msg": err.Error()})
			return
		}
		query := c.Request.URL.Query()
		fmt.Printf("%v, %v", reflect.TypeOf(query), query)
		data := Data{
			Message: Message{
				message.Channel,
				message.Action,
			},
			Query: query,
		}
		fmt.Printf("%v", data)
		c.JSON(200, gin.H{"channel": message.Channel, "action": message.Action})
	})
	return router
}
