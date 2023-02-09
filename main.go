package main

import (
	"qsr-mock-server/router"
	"qsr-mock-server/system_init"
)

func main() {
	system_init.Init()
	engine := router.Router()
	// 启动web服务，开启Mock监听
	engine.Run(":8787")
}
