package database

import (
	errcode "OneDisk/definition/err_code"
	"OneDisk/lib/log"
	"OneDisk/lib/random"
	timeutils "OneDisk/lib/utils/time"
	apiconstinvitecode "OneDisk/module/server/api/const/invitecode"
	"encoding/json"
	"go.uber.org/zap"
)

// InviteCodeCreateAndSaveForRegister
// 创建并保存一个用于注册账号的邀请码
func InviteCodeCreateAndSaveForRegister(
	fromUserID int64,
) (*InviteCode, OperationResult) {
	// 1、实例化激活码
	currentTimestamp := timeutils.Timestamp()
	extraModel := apiconstinvitecode.ExtraForRegister{
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
		ExpiredTime: currentTimestamp + apiconstinvitecode.TimeInviteCodeForRegisterDuration,
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
