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

// StorageFind
// 根据存储 ID 查找存储
func StorageFind(storageID int64) (*Storage, OperationResult) {
	var storage Storage
	findResult := database.First(&storage, storageID)
	if findResult.Error != nil {
		return nil, OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "Error occurred while finding storage at StorageFind()",
		}
	}
	return &storage, OperationResult{Code: errcode.OK}
}

// StorageListByCreatUserID
// 根据创建用户 ID 查找存储列表
func StorageListByCreatUserID(userID int64) ([]Storage, OperationResult) {
	// 1、查找用户是否合法
	queryUser, result := UserFindUserActiveStatus(userID)
	if result.Code != errcode.OK {
		return nil, result
	}
	if queryUser == nil {
		return nil, OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "User not found",
		}
	}
	// 2、查找存储列表
	var storageList []Storage
	queryStorageError := database.Where(&Storage{CreateUserID: userID}).Find(&storageList).Error
	if queryStorageError != nil {
		return nil, OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "Error occurred while finding storage at StorageListByCreatUserID()",
		}
	}
	return storageList, OperationResult{Code: errcode.OK}
}
