package database

import (
	errcode "OneDisk/definition/err_code"
	"OneDisk/lib/log"
	"OneDisk/lib/random"
	timeutils "OneDisk/lib/utils/time"
	apiconstinvitecode2 "OneDisk/server/api/const/invitecode"
	"encoding/json"
	"go.uber.org/zap"
)

// InviteCodeFindByID
// 根据 邀请码 ID 查找邀请码。找不到则返回 nil，但是 Code 为 OK
func InviteCodeFindByID(inviteCodeID int64) (*InviteCode, OperationResult) {
	var queryInviteCodes []InviteCode
	queryResult := database.Where(&InviteCode{ID: inviteCodeID}).Find(&queryInviteCodes)
	if queryResult.Error != nil {
		log.Info(tag, "Error occurred while finding invite code by ID in InviteCodeFindByID()", zap.Error(queryResult.Error))
		return nil, OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "Error occurred while finding invite code by ID",
			Error:   queryResult.Error,
		}
	}
	if len(queryInviteCodes) == 0 {
		return nil, OperationResult{Code: errcode.OK}
	}
	return &queryInviteCodes[0], OperationResult{Code: errcode.OK}
}

// InviteCodeCreateAndSaveForRegister
// 创建并保存一个用于注册账号的邀请码
func InviteCodeCreateAndSaveForRegister(fromUserID int64) (*InviteCode, OperationResult) {
	// 1、实例化激活码
	currentTimestamp := timeutils.Timestamp()
	extraModel := apiconstinvitecode2.ExtraForRegister{
		UserType: ValueUserTypeNormal,
	}
	extraModelJsonStr, err := json.Marshal(extraModel)
	if err != nil {
		log.Info(tag, "Error occurred while marshaling extra model in InviteCodeCreateAndSaveForRegister()", zap.Error(err))
		return nil, OperationResult{
			Code:    errcode.JSONConvert,
			Message: "Error occurred while marshaling extra model",
			Error:   err,
		}
	}
	insertInviteCode := InviteCode{
		FromUserID:  fromUserID,
		CreateTime:  currentTimestamp,
		ExpiredTime: currentTimestamp + apiconstinvitecode2.TimeInviteCodeForRegisterDuration,
		Usage:       ValueInviteCodeUsageRegister,
		Status:      ValueInviteCodeStatusNotUse,
		Code:        random.String(32),
		Extra:       string(extraModelJsonStr),
	}
	// 2、保存到数据库
	result := database.Create(&insertInviteCode)
	if result.Error != nil {
		log.Info(tag, "Error occurred while saving invite code in InviteCodeCreateAndSaveForRegister()", zap.Error(result.Error))
		return nil, OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "Error occurred while saving invite code",
			Error:   result.Error,
		}
	}
	return &insertInviteCode, OperationResult{Code: errcode.OK}
}

// InviteCodeInvalid
// 将邀请码标记为失效
func InviteCodeInvalid(userID int64, inviteCodeID int64) OperationResult {
	// 1、查找邀请码
	queryInviteCode, result := InviteCodeFindByID(inviteCodeID)
	if result.Code != errcode.OK {
		return result
	}
	// 2、检查邀请码是否存在
	if queryInviteCode == nil {
		return OperationResult{Code: errcode.DatabaseNotFound}
	}
	// 3、检查邀请码是否为自己的邀请码
	if queryInviteCode.FromUserID != userID {
		return OperationResult{Code: errcode.DatabaseNotFound}
	}
	// 4、检查邀请码是否已经失效
	if queryInviteCode.Status == ValueInviteCodeStatusInvalid {
		return OperationResult{Code: errcode.OK}
	}
	// 5、标记邀请码为失效并更新到数据库
	queryInviteCode.Status = ValueInviteCodeStatusInvalid
	saveResult := database.Save(queryInviteCode)
	if saveResult.Error != nil {
		log.Info(tag, "Error occurred while saving invite code in InviteCodeInvalid()", zap.Error(saveResult.Error))
		return OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "Error occurred while saving invite code",
			Error:   saveResult.Error,
		}
	}
	return OperationResult{Code: errcode.OK}
}
