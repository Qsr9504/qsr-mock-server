package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"qsr-mock-server/docs"
	"qsr-mock-server/service"
)

func Router() *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Any("/mock/*realUrl", service.ProxyHandler) // 获取所有拦截，进行自定义处理
	dbsRouter := r.Group("/mockRule")
	{

		dbsRouter.POST("/addRule", service.RuleAdd)                // 根据 ownerName 和 url 来进行存储
		dbsRouter.POST("/deleteRulesById", service.DeleteRuleById) // 根据 ID 来进行删除某一个指定的规则
		dbsRouter.POST("/updateRulesById", service.UpdateRuleById) // 根据 ID 来进行更新
		dbsRouter.GET("/getAllRules", service.GetAll)              // 查询所有的规则，无用户差异的 TODO：用户级别隔离
		dbsRouter.GET("/getRuleDetailById", service.GetRuleById)   // 根据 ID 获取某一个规则详情
		dbsRouter.GET("/getRuleByCons", service.GetAllByCons)      // 根据筛选条件进行查询
	}

	// 长链接Mock服务
	// ws://localhost:8787/spider-service/v1/channels
	spiderService := r.Group("/spider-service")
	{
		// v1 版本
		v1 := spiderService.Group("/v1")
		{
			v1.GET("/channels", service.NewClientConn)
		}
	}
	return r
}
