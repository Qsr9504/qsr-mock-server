package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

// GetHostAndPathFromCtx 根据 上下文信息 获取 remoteHost, remotePath
func getHostAndPathFromCtx(ctx *gin.Context) (string, string) {
	// 不需要判断是否包含 http:// 因为，客户端会直接报错 Error: Invalid protocol ，不允许这样进行请求
	// 如：localhost:8787/mock/localhost:8081/ping?hello=123&nihao=789&qsr=qqqqqqq 是错误的

	var (
		protocolScheme string
		remoteHost     string
	)

	if "HTTP/1.1" == ctx.Request.Proto {
		protocolScheme = "http://"
	} else {
		protocolScheme = "https://"
	}

	// 将realUrl进行拆分，获取Scheme和host , realUrl 已经过滤了 url 的参数信息
	strs := strings.Split(ctx.Param("realUrl"), "/")
	remotePathWithOutParams := ""
	for index, str := range strs {
		if index > 1 {
			remotePathWithOutParams += "/"
			remotePathWithOutParams += str
		}
		if index == 1 {
			// 截取 远程 路由地址
			remoteHost = str
		}
	}
	// 添加protocol cheme
	remoteHost = protocolScheme + remoteHost
	fmt.Println("---*", remoteHost, "*---")
	remotePathWithOutParams = strings.Split(remotePathWithOutParams, "?")[0] // 如果 url 中包含了参数，就进行过滤，这行原则上是没有用的
	return remoteHost, remotePathWithOutParams
}
