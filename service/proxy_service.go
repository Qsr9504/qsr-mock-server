package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
	"qsr-mock-server/models"
	"qsr-mock-server/utils"
)

// ProxyHandler
// @Summary	代理核心方法
// @Tags	代理模块
// @Router /mock/* [ANY]
func ProxyHandler(ctx *gin.Context) {
	remoteHost, remotePathWithOutParams := getHostAndPathFromCtx(ctx)
	// 从请求头中获取当前用户信息
	ownerName := ctx.Request.Header.Get("ownerName")
	fmt.Println("*** mock服务 ***\n 解析出来的 remoteHost=", remoteHost, " remotePathWithOutParams=", remotePathWithOutParams, "\n*** mock服务 ***\n")

	// 1. 根据 路由信息，从数据库中查询数据
	rulesModel := models.GetRulesByUserAndUrl(ownerName, remotePathWithOutParams)
	//fmt.Println("---checkSatisfyConditions(ctx, rulesModel.ProxyConditionJson)---", checkSatisfyConditions(ctx, rulesModel.ProxyConditionJson))
	// 2. 查询规则存在并且满足了拦截条件，才进行自定义mock数据返回
	if rulesModel.ProxyUrl != "" && checkSatisfyConditions(ctx, rulesModel.ProxyConditionJson) {
		fmt.Println("--- 拦截成功 ---")
		jsonMap, err := utils.JsonStringToMap(rulesModel.CustomizedResp) // 将自定义的返回值 字符串 捞出来转换为map
		if err != nil {
			fmt.Println("JSON 字符串 转 Map 出错，检查数据库中 字符串是否是后边有人改动过")
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, jsonMap) // 将自定义的返回值 设置为返回值
		return
	} else {
		// 拦截失败，直接转发，使用url中解析出来的远程ip地址转发
		var proxyUrl, _ = url.Parse(remoteHost)
		proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
		ctx.Request.URL.Path = remotePathWithOutParams // 修改上下文的路由地址
		ctx.Request.URL.Host = remoteHost              // 修改 上下文 中的 请求 主机地址，可以是host 也可以是 host:port
		proxy.ServeHTTP(ctx.Writer, ctx.Request)
		ctx.Abort()
	}
}

/*

func ProxyPostTest(ctx *gin.Context) {
	// 1. 解析 路由 和 远程地址
	// 2. 查询数据库，根据 路由 判断是否需要 拦截
	// 	2.1 不需要直接转发到指定地址，获取 返回值（错误信息） 并直接返回
	//	2.2 需要的话 进一步 校对 拦截条件
	// 3.校对成功，将预置的返回值进行返回
	// 4.校对失败，将直接转发到远程主机，获取 返回值（错误信息） 并直接返回
	ctx.JSON(200, gin.H{
		"message": "1234",
		"json":    "post",
	})
}

func ProxyTest(ctx *gin.Context) {
	// 1. 解析 路由 和 远程地址
	// 2. 查询数据库，根据 路由 判断是否需要 拦截
	// 	2.1 不需要直接转发到指定地址，获取 返回值（错误信息） 并直接返回
	//	2.2 需要的话 进一步 校对 拦截条件
	// 3.校对成功，将预置的返回值进行返回
	// 4.校对失败，将直接转发到远程主机，获取 返回值（错误信息） 并直接返回

	// 从body中获取参数
	value, exists := ctx.Get("qsr")
	fmt.Println("value=", value, " exists=", exists)

	value1, exists1 := ctx.GetPostForm("qsr")
	fmt.Println("value1=", value1, " exists1=", exists1)

	value2 := ctx.GetString("qsr")
	fmt.Println("value2=", value2)

	value3 := ctx.DefaultQuery("qsr", "999999")
	fmt.Println("value3=", value3)

	value4, exist4 := ctx.GetQuery("qsr")
	fmt.Println("value4=", value4, " exists4=", exist4)

	value5, exist5 := ctx.Get("qsr")
	fmt.Println("value5=", value5, " exist5=", exist5)

	value6 := ctx.Request.PostFormValue("qsr") // 可以动态根据请求头中的 application/x-www-form-urlencoded 或 multipart/form-data 动态进行获取 参数 --- 换句人话，除了json的 body中的数据都能拿
	fmt.Println("value6=", value6)

	data, _ := ioutil.ReadAll(ctx.Request.Body) // 从body-json中获取参数
	fmt.Println(string(data))

	ctx.JSON(200, gin.H{
		"message": "1234",
		"json":    "post",
	})
}

*/
