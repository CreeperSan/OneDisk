package config

import (
	"OneDisk/lib/definition"
	string2 "OneDisk/lib/format/formatstring"
	"OneDisk/lib/log"
	"gopkg.in/yaml.v2"
	"os"
)

func GetServer() Server {
	return cache.Server
}

func SetServer(server Server) error {
	// 更新到缓存中
	if len(server.Host) > 0 {
		cache.Server.Host = server.Host
	}
	if server.Port > 0 {
		cache.Server.Port = server.Port
	}

	// 转换为 yaml
	data, err := yaml.Marshal(cache)
	if err != nil {
		log.Error(tag, string2.String("Failed to marshal config.yaml when set server, server=", server))
		return err
	}

	// 写入文件
	err = os.WriteFile(definition.PathConfig, data, 0644)
	if err != nil {
		log.Error(tag, string2.String("Failed to write config.yaml when set server, server=", server))
		return err
	}

	return nil
}
