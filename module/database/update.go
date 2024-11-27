package database

import (
	"OneDisk/lib/definition"
	string2 "OneDisk/lib/format/formatstring"
	"OneDisk/lib/log"
	"database/sql"
	"go.uber.org/zap"
	"strconv"
)

func getDatabaseVersion(db *sql.DB) (int, error) {
	var version string
	err := db.QueryRow("SELECT sqlite_version()").Scan(&version)
	if err != nil {
		return -1, err
	}
	versionCode, err := strconv.Atoi(version)
	if err != nil {
		return -1, err
	}
	return versionCode, nil
}

func upgradeDatabase(db *sql.DB) error {
	// 获取数据库版本
	databaseVersion, err := getDatabaseVersion(db)
	if err != nil {
		log.Error(tag, "Failed to get database version", zap.Error(err))
	}

	log.Info(tag, string2.String("Current database version is %d", databaseVersion))

	if databaseVersion >= definition.VersionDatabaseLatest {
		log.Info(tag, "Database is up to date")
		return nil
	} else if databaseVersion < definition.VersionDatabaseInitialize {
		// 数据库初始化
		db.Exec("CREATE TABLE IF NOT EXISTS user (" +
			"id INTEGER PRIMARY KEY AUTOINCREMENT" + "," +
			"name TEXT NOT NULL" + "," +
			"email TEXT NOT NULL DEFAULT ''" + "," +
			"avatar TEXT NOT NULL DEFAULT ''" +
			")")
		db.Exec("INSERT INTO user (name) VALUES ('admin')")
	}

	return nil
}
