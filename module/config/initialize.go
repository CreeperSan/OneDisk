package config

import (
	"OneDisk/lib/definition"
	"OneDisk/lib/log"
	"OneDisk/lib/utils/file"
	"errors"
)

const tag = "Config"

func Initialize() error {
	if !fileutils.Exists(definition.PathConfig) {
		log.Info(tag, "Config.yaml not found, creating...")
		if !fileutils.CreateFile(definition.PathConfig) {
			log.Error(tag, "Failed to create config.yaml")
			return errors.New("failed to create config.yaml")
		} else {
			log.Info(tag, "Config.yaml created successfully")
		}
	} else {
		log.Info(tag, "Config.yaml founded")
	}
	return nil
}
