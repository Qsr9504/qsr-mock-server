package service

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"qsr-mock-server/models"
	"qsr-mock-server/utils"
)

// 判断当前 路由 的拦截条件，是否 拦截成功，true 满足了拦截条件，即拦截成功，需要走mock的预置消息
// -- 如果报错，就是拦截失败，返回 false
func checkSatisfyConditions(ctx *gin.Context, proxyConditionJson string) bool {
	fmt.Println("---proxyConditionJson---", proxyConditionJson)
	if proxyConditionJson == "" {
		// 如果条件是""，说明拦截失败，条件直接不满足
		return false
	}
	// 将 proxyConditionJson 解析 为 ProxyCondition 对象集合
	conditions, err := models.JsonToProxyConditions(proxyConditionJson)
	fmt.Println("---conditions---", conditions)
	if err != nil {
		fmt.Println("数据库中获取出来的 条件json 转 []ProxyCondition 失败", err)
		return false
	}

	conditionPass := false // 条件检查，只要有一个不满足，即false，就是拦截失败
	for _, condition := range conditions {
		switch condition.InterceptSite {
		case models.Header:
			// 这个拦截条件是从header中拦截
			conditionPass = headerFilter(ctx, condition)
			fmt.Println("---Header---", conditionPass)
		case models.Uri:
			// 拦截条件是从URI中获取
			conditionPass = uriFilter(ctx, condition)
		case models.Body:
			// 拦截条件是从 Body 中获取
			conditionPass = bodyFilter(ctx, condition)
			fmt.Println("---Body---", conditionPass)
		}
		if !conditionPass {
			// 只要有一次是 false，整个判断结束，拦截失败，返回 false
			return false
		}
	}
	return true
}

func headerFilter(ctx *gin.Context, proxyCondition models.ProxyCondition) bool {
	if proxyCondition.IsContain == (ctx.GetHeader(proxyCondition.InterceptKey) == proxyCondition.InterceptValue) {
		return true
	}
	return false
}

// uriFilter url 中的参数获取，判断是否存在
func uriFilter(ctx *gin.Context, proxyCondition models.ProxyCondition) bool {
	// 从请求 中获取 Query ，如果key
	if proxyCondition.IsContain == (ctx.Query(proxyCondition.InterceptKey) == proxyCondition.InterceptValue) {
		return true
	}
	return false
}

// bodyFilter 的 body 参数中进行获取判断，是否又满足筛选条件的key
func bodyFilter(ctx *gin.Context, proxyCondition models.ProxyCondition) bool {
	if ctx.GetHeader("Content-Type") == "application/json" {
		bodyByte, _ := ioutil.ReadAll(ctx.Request.Body)
		// 将流中信息读取后需要重新放进去,要不然读取完，缓冲区就空了
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyByte))
		keyAndValueExist := utils.JsonStrContain(string(bodyByte), proxyCondition.InterceptKey, proxyCondition.InterceptValue)
		return keyAndValueExist == proxyCondition.IsContain
	} else {
		// body-form 类型的请求，这样获取
		keyAndValueExist := ctx.Request.PostFormValue(proxyCondition.InterceptKey) == proxyCondition.InterceptValue
		return keyAndValueExist == proxyCondition.IsContain
	}
}

// GetJSONBodyFromCtx 从上下文的 body-json 中获取数据
func GetJSONBodyFromCtx(ctx *gin.Context) (map[string]interface{}, error) {
	data, _ := ioutil.ReadAll(ctx.Request.Body) // 从body-json中获取参数
	// 将流中信息读取后需要重新放进去,再重新写回请求体body中，ioutil.ReadAll会清空c.Request.Body中的数据
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	return utils.JsonStringToMap(string(data))
}
