package config

import (
	"OneDisk/definition"
	"OneDisk/lib/log"
	"OneDisk/lib/utils/file"
	"errors"
	"gopkg.in/yaml.v2"
	"os"
)

const tag = "Config"

func Initialize() error {
	log.Info(tag, "Config initializing...")

	const pathConfig = definition.PathConfig
	// 文件检查
	if !fileutils.Exists(pathConfig) {
		log.Info(tag, "Config.yaml not found, creating...")
		if !fileutils.CreateFile(pathConfig) {
			log.Error(tag, "Failed to create config.yaml")
			return errors.New("failed to create config.yaml")
		} else {
			log.Info(tag, "Config.yaml created successfully")
		}
	} else {
		log.Info(tag, "Config.yaml founded")
	}
	// 读取配置文件
	configData, err := os.ReadFile(pathConfig)
	if err != nil {
		log.Error(tag, "Failed to read config.yaml")
		return err
	}
	// 解析配置文件
	err = yaml.Unmarshal(configData, &cache)
	if err != nil {
		log.Error(tag, "Failed to unmarshal config.yaml")
		return err
	}

	log.Info(tag, "Config initialized successfully")
	return nil
}
