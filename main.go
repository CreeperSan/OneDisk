package main

import (
	"OneDisk/definition"
	"OneDisk/lib/lifecycle"
	"OneDisk/lib/log"
	"OneDisk/module/config"
	"OneDisk/module/database"
	"OneDisk/server"
)

func main() {
	// 日志初始化
	log.Initialize() // 先初始化日志
	log.AppStart()   // 打印应用启动日志

	// 配置文件初始化
	err := config.Initialize()
	if err != nil {
		lifecycle.Exit(definition.ExitCodeConfigInitialize)
		return
	}

	// 数据库配置初始化
	err = database.Initialize()
	if err != nil {
		lifecycle.Exit(definition.ExitCodeDatabaseInitialize)
		return
	}

	// 服务器配置初始化
	err = server.Initialize()
	if err != nil {
		lifecycle.Exit(definition.ExitCodeServerInitialize)
		return
	}

	// 启动服务器
	err = server.StartServer()
	if err != nil {
		lifecycle.Exit(definition.ExitCodeServerStart)
		return
	}
}
