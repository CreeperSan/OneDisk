package database

import (
	"OneDisk/lib/definition"
	"OneDisk/lib/format/formatstring"
	string2 "OneDisk/lib/format/formatstring"
	"OneDisk/lib/log"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
)

func getDatabaseVersion(db *gorm.DB) (int, error) {
	var version string
	result := db.Raw("PRAGMA user_version").Scan(&version)
	if result.Error != nil {
		return -1, result.Error
	}
	versionCode, err := strconv.Atoi(version)
	if err != nil {
		return -1, err
	}
	return versionCode, nil
}

func setDatabaseVersion(db *gorm.DB, version int) error {
	return db.Exec(formatstring.String("PRAGMA user_version = %d", version)).Error
}

func checkAndUpgradeVersion(db *gorm.DB) error {
	// 获取数据库版本
	databaseVersion, err := getDatabaseVersion(db)
	if err != nil {
		log.Error(tag, "Failed to get database version", zap.Error(err))
	}

	log.Info(tag, string2.String("Current database version is %d", databaseVersion))

	for databaseVersion < definition.VersionDatabaseLatest {
		latestVersion, err := upgradeDatabase(db, databaseVersion)
		if err != nil {
			log.Error(tag, "Failed to upgrade database", zap.Error(err))
			return err
		}
		err = setDatabaseVersion(db, latestVersion)
		if err != nil {
			log.Error(tag, "Failed to set database version", zap.Error(err))
			return err
		}
		log.Info(tag, string2.String("Database upgraded from version %d to version %d", databaseVersion, latestVersion))
		databaseVersion = latestVersion
	}

	return nil
}

func upgradeDatabase(db *gorm.DB, currentVersion int) (int, error) {
	if currentVersion < definition.VersionDatabaseInitialize {
		/* 数据库初版初始化 */
		// 创建用户表
		db.Exec("CREATE TABLE IF NOT EXISTS " + tableUser + " (" +
			columnUserID + " INTEGER PRIMARY KEY AUTOINCREMENT," +
			columnUserUsername + " VARCHAR(64) NOT NULL UNIQUE," +
			columnUserPassword + " VARCHAR(128) NOT NULL," +
			columnUserEmail + " VARCHAR(128) UNIQUE," +
			columnUserNickname + " VARCHAR(128) NOT NULL," +
			columnUserAvatar + " VARCHAR(256)," +
			columnUserPhone + " VARCHAR(32) UNIQUE," +
			columnUserCreateTime + " INTEGER NOT NULL," +
			columnUserType + " INTEGER NOT NULL DEFAULT " + strconv.Itoa(valueUserTypeGuest) +
			")")
		// 创建用户令牌表
		db.Exec("CREATE TABLE IF NOT EXISTS " + tableUserToken + " (" +
			columnUserTokenID + " INTEGER PRIMARY KEY AUTOINCREMENT," +
			columnUserTokenUserID + " INTEGER NOT NULL," +
			columnUserTokenToken + " VARCHAR(128) NOT NULL," +
			columnUserTokenPlatform + " INTEGER NOT NULL," +
			columnUserTokenMachineCode + " VARCHAR(32) NOT NULL," +
			columnUserTokenMachineName + " VARCHAR(32) NOT NULL," +
			columnUserTokenSecretKey + " VARCHAR(32) NOT NULL," +
			columnUserTokenCreateTime + " INTEGER NOT NULL," +
			columnUserTokenExpireTime + " INTEGER NOT NULL," +
			"FOREIGN KEY (" + columnUserTokenUserID + ") REFERENCES " + tableUser + "(" + columnUserID + ")" +
			")")
		return definition.VersionDatabaseInitialize, nil
	}
	return currentVersion, nil
}
