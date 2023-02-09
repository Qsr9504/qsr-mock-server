package models

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"qsr-mock-server/common"
)

const (
	Header = iota
	Body
	Uri
)

type ProxyCondition struct {
	InterceptSite  int    // 拦截的目标位置，Header、Body、Uri 这几种方式
	IsContain      bool   // 包含、不包含
	InterceptKey   string // 拦截的 目标 key
	InterceptValue string // 拦截的 目标 value
}

type ProxyRulesModel struct {
	gorm.Model
	OwnerName          string           // 所属的用户
	ProxyMethod        string           // 需要代理的方法
	ProxyHost          string           // 需要代理的主机地址 ip + port
	ProxyUrl           string           // 需要代理的路由地址，不包含请求参数 和 http:// 或 https://
	ProxyConditionJson string           // 代理条件，json字符串，可以传入多个条件
	ProxyConditions    []ProxyCondition `gorm:"-"` // 代理的条件集合, gorm 忽略此字段
	CustomizedResp     string           // 定制的返回值
}

func (table *ProxyRulesModel) TableName() string {
	return common.ProxyRulesModelTableName
}

// GetDemoJsonReq 获取请求字符串 样例
func (pRuleModel *ProxyRulesModel) GetDemoJsonReq() (string, error) {
	bytes, err := json.Marshal(pRuleModel)
	return string(bytes), err
}

// GetRulesByUserAndUrl 根据   规则所主 和 url   同时判断 规则是否存在
func GetRulesByUserAndUrl(owner, url string) ProxyRulesModel {
	rulesModel := ProxyRulesModel{}
	common.DB.Where("owner_name = ? and proxy_url = ?", owner, url).First(&rulesModel)
	return rulesModel
}

// CreateRules 创建一个拦截规则
func CreateRules(proxyRule *ProxyRulesModel) (*gorm.DB, *ProxyRulesModel) {
	// 插入之前先进行检查，如果数据已经存在，就不再进行插入操作
	rule := GetRulesByUserAndUrl(proxyRule.OwnerName, proxyRule.ProxyUrl)
	fmt.Println(rule)
	if rule.ID != 0 {
		return nil, &rule
	}
	tx := common.DB.Create(proxyRule)
	return tx, nil
}

// DeleteRulesById 根据id 删除一个拦截规则,逻辑删
func DeleteRulesById(id uint) *gorm.DB {
	rulesById := GetRulesById(id)
	return common.DB.Delete(&rulesById)
}

// GetRulesById 通过 id 获取一个拦截规则
func GetRulesById(id uint) ProxyRulesModel {
	rulesModel := ProxyRulesModel{}
	common.DB.Where("id = ?", id).First(&rulesModel)
	return rulesModel
}

// GetRulesByCons 按照条件进行查询
//func GetRulesByCons(cons map[string]interface{}) []ProxyRulesModel {
func GetRulesByCons(rule *ProxyRulesModel) []ProxyRulesModel {
	var rulesModels []ProxyRulesModel
	common.DB.Where(rule).Find(&rulesModels)
	return rulesModels
}

// GetAllRules 获取所有规则
// @Resp  rulesModels: 数据库中捞到的数据  result.RowsAffected： 影响行数
// TODO: 加入分页查询
func GetAllRules() ([]ProxyRulesModel, int64) {
	var rulesModels []ProxyRulesModel
	result := common.DB.Find(&rulesModels)
	// 如果没有问题的话，result 就是 nil
	fmt.Println(result)
	if result.Error != nil {
		return nil, 0
	} else {
		return rulesModels, result.RowsAffected
	}
}

// UpdateRuleById 通过规则id 更新一个拦截规则
func UpdateRuleById(proxyRule ProxyRulesModel) *gorm.DB {
	db := common.DB.Model(&proxyRule).Updates(proxyRule).Where("id=?", proxyRule.ID)
	return db
}

// JsonToProxyConditions 将一个 ProxyConditionJson 转换为 []ProxyCondition 类型
func JsonToProxyConditions(proxyConditionJson string) ([]ProxyCondition, error) {
	var proxyConditions []ProxyCondition
	err := json.Unmarshal([]byte(proxyConditionJson), &proxyConditions)
	if err != nil {
		return nil, err
	}
	return proxyConditions, nil
}
