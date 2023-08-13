package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"learn_zinx/Cobra.mayfly/internal/common/api"
	dbapp "learn_zinx/Cobra.mayfly/internal/db/application"
	machineapp "learn_zinx/Cobra.mayfly/internal/machine/application"
	projectapp "learn_zinx/Cobra.mayfly/internal/project/application"
	redisapp "learn_zinx/Cobra.mayfly/internal/redis/application"
	"learn_zinx/Cobra.mayfly/pkg/ctx"
	"reflect"
)

// 原来的初始化方法
func InitIndexRouter(router *gin.RouterGroup) {
	index := router.Group("common/index")

	i := &api.Index{
		ProjectApp: projectapp.GetProjectApp(),
		MachineApp: machineapp.GetMachineApp(),
		DbApp:      dbapp.GetDbApp(),
		RedisApp:   redisapp.GetRedisApp(),
		Val:        1,
	}

	vi := reflect.ValueOf(i)
	//fmt.Printf("name:'%v' kind:'%v'\n", typeOfi.Name(), typeOfi.Kind())

	//typeOfv := typeOfi.Elem()
	//vt := reflect.ValueOf(typeOfv)
	vn := vi.MethodByName("Test")
	// 如果存在情况下 is OK? false false true
	if vn.IsValid() && !vn.IsNil() { // 判断获得的方法是否可用
		result := vn.Call([]reflect.Value{reflect.ValueOf(1)})
		fmt.Println("结果为:", result[0].Int())
	}

	//fmt.Println("isOk?", vn)

	{ // 形成一个作用域，里面会覆盖外边
		// 首页基本信息统计
		index.GET("count", func(g *gin.Context) {
			ctx.NewReqCtxWithGin(g).
				Handle(i.Count)
		})
	}
}
