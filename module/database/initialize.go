package database

import (
	"OneDisk/lib/definition"
	"OneDisk/lib/format/formatstring"
	"OneDisk/lib/input"
	"OneDisk/lib/log"
	"OneDisk/lib/random"
	"OneDisk/module/config"
	"fmt"
	sqliteEncrypt "github.com/hinha/gorm-sqlite-cipher"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var tag = "Database"

var database *gorm.DB

func Initialize() error {
	log.Info(tag, "Database initializing...")

	// 检查数据库配置
	configDatabase := config.GetDatabase()
	if len(configDatabase.Type) <= 0 {
		// 未配置数据库，进行配置
		log.Info(tag, "Database not found, creating ...")
		// 配置类型（目前仅支持sqlite）
		configDatabase.Type = definition.DatabaseSqlite
		// 配置路径（目前暂不支持自定义）
		configDatabase.Path = definition.PathDatabase
		// 配置密码
		fmt.Println("Please enter your database password (leave empty will pick random password):")
		inputPassword := input.ReadString()
		if len(inputPassword) <= 0 {
			inputPassword = random.Password(20)
		}
		configDatabase.Password = inputPassword

		err := config.SetDatabase(configDatabase)
		if err != nil {
			log.Error(tag, "Failed to set database")
			return err
		}
	}

	// 尝试连接数据库
	log.Info(tag, "Opening database...")
	sqlOpenDatabase := formatstring.String("%s?_pragma_key=%s&_pragma_cipher_page_size=4096", configDatabase.Path, configDatabase.Password)
	db, err := gorm.Open(sqliteEncrypt.Open(sqlOpenDatabase), &gorm.Config{})
	if err != nil {
		log.Error(tag, "Failed to open database", zap.Error(err))
		return err
	}

	// 数据库升级gorm
	//log.Info(tag, "Checking database update...")
	//err = upgradeDatabase(db)
	//if err != nil {
	//	log.Error(tag, "Failed to upgrade database", zap.Error(err))
	//	return err
	//}

	database = db

	log.Info(tag, "Database initialized successfully")
	return nil
}
