# ginSwagger
> 官网：https://pkg.go.dev/github.com/swaggo/gin-swagger#section-readme

案例：
```
// @BasePath /api/v1
// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func Helloworld(g *gin.Context)  {
    g.JSON(http.StatusOK,"helloworld")  
}
```

1. 路由配置
```
docs.SwaggerInfo.BasePath = ""
engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

```

2. 