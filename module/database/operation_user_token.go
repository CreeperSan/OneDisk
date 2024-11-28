package database

import (
	errcode "OneDisk/definition/err_code"
	"OneDisk/lib/format/formatstring"
	"OneDisk/lib/log"
	timeutils "OneDisk/lib/utils/time"
	"go.uber.org/zap"
)

// UserTokenValidation
// 用户 Token 验证
func UserTokenValidation(
	userID int64,
	token string,
	machineCode string,
	machineName string,
	platform int,
) (*User, *UserToken, OperationResult) {
	// 查询用户是否被封禁或者注销
	var queryUsers []User
	queryResult := database.Where(&User{
		ID: userID,
	}).Find(&queryUsers)
	if queryResult.Error != nil {
		return nil, nil, OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "Error occurred while querying user",
			Error:   queryResult.Error,
		}
	}
	if len(queryUsers) <= 0 {
		return nil, nil, OperationResult{
			Code:    errcode.UserNotExist,
			Message: "User not exist",
		}
	}
	queryUser := queryUsers[0]
	if queryUser.Status == valueUserStatusForbidden {
		return nil, nil, OperationResult{
			Code:    errcode.UserForbidden,
			Message: "User is forbidden",
		}
	}

	// 查询 Token 是否合法
	var queryUserTokens []UserToken
	queryResult = database.Where(&UserToken{
		UserID:      userID,
		Token:       token,
		MachineCode: machineCode,
		MachineName: machineName,
		Platform:    platform,
	}).Find(&queryUserTokens)
	if queryResult.Error != nil {
		return nil, nil, OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "Error occurred while querying user token",
			Error:   queryResult.Error,
		}
	}
	if len(queryUserTokens) <= 0 {
		return nil, nil, OperationResult{
			Code:    errcode.AuthTokenInvalid,
			Message: "No token was found",
		}
	}
	queryUserToken := queryUserTokens[0]
	timestampCurrent := timeutils.Timestamp()
	if timestampCurrent > queryUserToken.ValidTime+queryUserToken.Duration {
		// 删除 Token
		queryResult = database.Delete(&queryUserToken)
		if queryResult.Error != nil {
			log.Warming(tag, formatstring.String("Fail to remove expired token: token=%s id=%d userID=%d", queryUserToken.Token, queryUserToken.ID, queryUserToken.UserID), zap.Error(queryResult.Error))
		}
		return nil, nil, OperationResult{
			Code:    errcode.AuthTokenExpired,
			Message: "Token expired",
		}
	}
	// 更新 Token 校验时间
	queryUserToken.ValidTime = timestampCurrent
	queryResult = database.Save(&queryUserToken)
	if queryResult.Error != nil {
		log.Warming(tag, formatstring.String("Failed to update token valid time.token=%s userID=%d", queryUserToken.Token, userID), zap.Error(queryResult.Error))
	}

	return &queryUser, &queryUserToken, OperationResult{Code: errcode.OK}
}
