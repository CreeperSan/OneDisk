package config

import (
	definition2 "OneDisk/definition"
	string2 "OneDisk/lib/format/formatstring"
	"OneDisk/lib/log"
	"errors"
	"gopkg.in/yaml.v2"
	"os"
)

func GetDatabase() Database {
	return cache.Database
}

func SetDatabase(database Database) error {
	// 更新到缓存中
	if database.Type == definition2.DatabaseSqlite { // 目前仅支持 sqlite
		if len(database.Path) <= 0 {
			database.Path = definition2.PathDatabase
		}
		cache.Database = database
	} else {
		return errors.New("unsupported database type")
	}

	// 转换为 yaml
	data, err := yaml.Marshal(cache)
	if err != nil {
		log.Error(tag, string2.String("Failed to marshal config.yaml when set database, database=", database))
		return err
	}

	// 写入文件
	err = os.WriteFile(definition2.PathConfig, data, 0644)
	if err != nil {
		log.Error(tag, string2.String("Failed to write config.yaml when set database, database=", database))
		return err
	}

	return nil
}
