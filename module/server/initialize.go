package server

import (
	string2 "OneDisk/lib/format/formatstring"
	"OneDisk/lib/input"
	"OneDisk/lib/log"
	"OneDisk/module/config"
	apiinvitecode "OneDisk/module/server/api/invitecode"
	apimiddleware "OneDisk/module/server/api/middleware"
	apiuser "OneDisk/module/server/api/user"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
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
		log.Info(tag, string2.String("You have entered server host name: %s", inputHostName))
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
		log.Info(tag, string2.String("You have entered server port: %d", inputPort))
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

func StartServer() error {
	server := gin.Default()

	server.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// 公共中间件
	server.Use(apimiddleware.HeaderConvert())

	// 模块接口注册
	apiuser.Register(server)
	apiinvitecode.Register(server)

	err := server.Run()
	if err != nil {
		log.Error(tag, "Failed to start server", zap.Error(err))
		return err
	}

	return nil
}
