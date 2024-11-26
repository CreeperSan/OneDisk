package server

import (
	"OneDisk/lib/format"
	"OneDisk/lib/input"
	"OneDisk/lib/log"
	"OneDisk/module/config"
	"fmt"
)

const tag = "Server"

func Initialize() error {
	//server := config.GetServer()

	// 读取检查配置
	configServer := config.GetServer()

	// 检查主机名
	if len(configServer.Host) <= 0 {
		log.Info(tag, "Server host name not define, waiting for enter ...")
		fmt.Println("Please enter server host name:")
		inputHostName := input.ReadString()
		log.Info(tag, format.String("You have entered server host name: %s", inputHostName))
		if len(inputHostName) <= 0 {
			log.Error(tag, "The server host name you have entered is invalid")
			return fmt.Errorf("server host name can not be empty")
		}
		configServer.Host = inputHostName
		err := config.SetServer(configServer)
		if err != nil {
			log.Error(tag, "Failed to set server host name")
			return err
		}
	}

	// 检查端口号
	if configServer.Port <= 0 {
		log.Info(tag, "Server port not define, waiting for enter ...")
		fmt.Println("Please enter server port:")
		inputPort := input.ReadInt()
		log.Info(tag, format.String("You have entered server port: %d", inputPort))
		if inputPort <= 0 {
			log.Error(tag, "The server port you have entered is invalid")
			return fmt.Errorf("server port can not be empty")
		}
		configServer.Port = inputPort
		err := config.SetServer(configServer)
		if err != nil {
			log.Error(tag, "Failed to set server port")
			return err
		}
	}

	return nil
}
