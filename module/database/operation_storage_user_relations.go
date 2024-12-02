package database

import (
	errcode "OneDisk/def/err_code"
	timeutils "OneDisk/lib/utils/time"
)

// StorageUserRelationFindByUserIDAndStorageID
// 根据用户 ID 和存储 ID 查找用户关系
func StorageUserRelationFindByUserIDAndStorageID(userID int64, storageID int64) (*StorageUserRelation, OperationResult) {
	var queryRelations []StorageUserRelation
	result := database.Where(&StorageUserRelation{UserID: userID, StorageID: storageID}).Find(&queryRelations)
	if result.Error != nil {
		return nil, OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "Error occurred while finding storage user relation at StorageUserRelationFindByUserIDAndStorageID()",
		}
	}
	if len(queryRelations) == 0 {
		return nil, OperationResult{Code: errcode.OK}
	}
	queryRelation := queryRelations[0]
	return &queryRelation, OperationResult{Code: errcode.OK}
}

// StorageUserRelationCreateIfEmpty
// 如果用户关系不存在则创建
func StorageUserRelationCreateIfEmpty(userID int64, storageID int64) (*StorageUserRelation, OperationResult) {
	// 1、查询用户是否存在
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
	// 2、查询存储是否存在
	queryStorage, result := StorageFind(storageID)
	if result.Code != errcode.OK {
		return nil, result
	}
	if queryStorage == nil {
		return nil, OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "Storage not found",
		}
	}
	// 3、查询用户关系是否存在
	queryRelation, result := StorageUserRelationCreateIfEmpty(userID, storageID)
	if result.Code != errcode.OK {
		return nil, result
	}
	if queryRelation != nil {
		// 3.1、已经存在则直接返回
		return queryRelation, OperationResult{Code: errcode.OK}
	}
	// 4、创建用户关系
	currentTimestamp := timeutils.Timestamp()
	insertStorageUserRelation := StorageUserRelation{
		UserID:     userID,
		StorageID:  storageID,
		CreateTime: currentTimestamp,
	}
	insertResult := database.Create(&insertStorageUserRelation)
	if insertResult.Error != nil {
		return nil, OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "Error occurred while saving storage user relation at StorageUserRelationCreateIfEmpty()",
		}
	}
	return &insertStorageUserRelation, OperationResult{Code: errcode.OK}
}

// StorageUserRelationList
// 获取用户的存储关系列表
func StorageUserRelationList(userID int64) ([]Storage, OperationResult) {
	// 1、查找用户是否存在
	queryUser, result := UserFindUser(userID)
	if result.Code != errcode.OK {
		return nil, result
	}
	if queryUser == nil {
		return nil, OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "User not found",
		}
	}
	// 2、获取对应关系
	var queryStorages []Storage
	queryStoragesError := database.Table(tableStorageUserRelation).
		Select(tableStorage+".*").
		Joins("LEFT JOIN "+tableStorage+" ON "+tableStorage+"."+columnUserID+" = "+tableStorageUserRelation+"."+columnStorageUserRelationUserID).
		Where(tableStorageUserRelation+"."+columnStorageUserRelationUserID+" = ?", userID).
		Scan(&queryStorages).Error
	if queryStoragesError != nil {
		return nil, OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "Error occurred while finding storage user relation at StorageUserRelationList()",
		}
	}
	return queryStorages, OperationResult{Code: errcode.OK}
}
