package database

import (
	errcode "OneDisk/def/err_code"
	timeutils "OneDisk/lib/utils/time"
)

// StorageCreateAndSaveForLocalPath
// 创建并保存一个本地路径的存储策略
func StorageCreateAndSaveForLocalPath(name string, avatar string, userID int64, storageType int, config string) (*Storage, OperationResult) {
	// 1、创建存储
	currentTimestamp := timeutils.Timestamp()
	storage := Storage{
		CreateUserID: userID,
		Name:         name,
		Avatar:       avatar,
		Type:         storageType,
		Config:       config,
		CreateTime:   currentTimestamp,
		UpdateTime:   currentTimestamp,
	}
	// 2、保存存储
	saveResult := database.Create(&storage)
	if saveResult.Error != nil {
		return nil, OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "Error occurred while saving storage at StorageCreateAndSaveForLocalPath()",
		}
	}
	return &storage, OperationResult{Code: errcode.OK}
}
