package common

import "gorm.io/gorm"

const APPName = "qsr-mock-server"

var (
	DB *gorm.DB
)

const (
	ProxyRulesModelTableName = "proxy_rules"
)
