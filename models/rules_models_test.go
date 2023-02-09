package models

import (
	"fmt"
	"testing"
	"time"
)

//type ProxyCondition struct {
//	InterceptSite  int    // 拦截的目标位置，Header、Body、Uri 这几种方式
//	IsContain      bool   // 包含、不包含
//	InterceptKey   string // 拦截的 目标 key
//	InterceptValue string // 拦截的 目标 value
//}
//
//type ProxyRulesModel struct {
//	gorm.Model
//	OwnerName          string           // 所属的用户
//	ProxyMethod        string           // 需要代理的方法
//	ProxyUrl           string           // 需要代理的路由地址
//	ProxyConditionJson string           // 代理条件，json字符串，可以传入多个条件
//	ProxyConditions    []ProxyCondition // 代理的条件集合
//	CustomizedResp     string           // 定制的返回值
//}

func TestGetDemoJsonReq(t *testing.T) {

	rulesModel := ProxyRulesModel{}
	rulesModel.CreatedAt = time.Now()
	rulesModel.UpdatedAt = time.Now()
	rulesModel.OwnerName = "demoName"
	rulesModel.ProxyMethod = "POST/ANY/GET"
	rulesModel.ProxyUrl = "http://localhost:8081/ping"
	rulesModel.CustomizedResp = "{\"message\":\"ok\"}"
	rulesModel.ProxyConditions = []ProxyCondition{}

	for i := 0; i < 2; i++ {
		pCon := ProxyCondition{
			InterceptSite:  Header,
			IsContain:      true,
			InterceptKey:   fmt.Sprint(i),
			InterceptValue: fmt.Sprint(i),
		}
		rulesModel.ProxyConditions = append(rulesModel.ProxyConditions, pCon)
	}
	req, err := rulesModel.GetDemoJsonReq()
	if err != nil {
		fmt.Println(err)
	}
	t.Log(req)
}
