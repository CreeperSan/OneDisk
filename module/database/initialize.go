package database

import (
	definition2 "OneDisk/definition"
	"OneDisk/lib/format/formatstring"
	"OneDisk/lib/input"
	"OneDisk/lib/log"
	"OneDisk/lib/random"
	timeutils "OneDisk/lib/utils/time"
	"OneDisk/module/config"
	"errors"
	"fmt"
	sqliteEncrypt "github.com/hinha/gorm-sqlite-cipher"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
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
		configDatabase.Type = definition2.DatabaseSqlite
		// 配置路径（目前暂不支持自定义）
		configDatabase.Path = definition2.PathDatabase
		// 配置密码
		fmt.Println("Please enter your database password:")
		inputPassword := input.ReadString()
		if len(inputPassword) <= 0 {
			fmt.Println("Are you sure to use no password? (enter y to confirm or any other key to generate a random password)")
			inputConfirmNoPassword := strings.ToLower(input.ReadString())
			if inputConfirmNoPassword != "y" { // 如果没有确认使用空密码，则随机生成32位密码
				inputPassword = random.Password(32)
			}
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
	log.Info(tag, "Checking database update...")
	err = checkAndUpgradeVersion(db)
	if err != nil {
		log.Error(tag, "Failed to upgrade database", zap.Error(err))
		return err
	}

	// 初始化管理员账户
	var queryAdminUser []User
	db.Where(formatstring.String("%s = ?", columnUserType), valueUserTypeAdmin).Find(&queryAdminUser)
	if queryAdminUser == nil || len(queryAdminUser) <= 0 {
		log.Info(tag, "No administrator account found, creating...")
		fmt.Println("You haven't created an administrator account yet.")
		fmt.Println("Please enter administrator username:")
		inputAdminUsername := input.ReadString()
		fmt.Println("Please enter administrator password:")
		inputAdminPassword := input.ReadString()
		tmpAdminUser := User{
			Username:   inputAdminUsername,
			Password:   formatstring.Password(inputAdminPassword),
			Nickname:   inputAdminUsername,
			CreateTime: timeutils.Timestamp(),
			Type:       valueUserTypeAdmin,
			Status:     valueUserStatusActive,
		}
		resultInsert := db.Create(&tmpAdminUser)
		if resultInsert.Error != nil {
			log.Error(tag, "Failed to create administrator account", zap.Error(resultInsert.Error))
			return resultInsert.Error
		}
		log.Info(tag, "Administrator account created successfully!", zap.String("username", inputAdminUsername))
	} else if queryAdminUser != nil && len(queryAdminUser) > 1 {
		log.Error(tag, "More than one administrator account found!")
		return errors.New("can not have more than one administrator account, please check the database")
	}

	database = db

	log.Info(tag, "Database initialized successfully")
	return nil
}
