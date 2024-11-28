package database

import (
	errcode "OneDisk/definition/err_code"
	"OneDisk/lib/format/formatstring"
	"OneDisk/lib/log"
	timeutils "OneDisk/lib/utils/time"
	apiconstuser "OneDisk/module/server/api/const/user"
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
	if timestampCurrent > queryUserToken.TokenExpireTime {
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
	queryUserToken.LastAccessTime = timestampCurrent
	queryResult = database.Save(&queryUserToken)
	if queryResult.Error != nil {
		log.Warming(tag, formatstring.String("Failed to update token access time.token=%s userID=%d", queryUserToken.Token, userID), zap.Error(queryResult.Error))
	}

	return &queryUser, &queryUserToken, OperationResult{Code: errcode.OK}
}

// UserTokenRefresh
// 刷新用户身份令牌
func UserTokenRefresh(
	userID int64,
	refreshToken string,
	platform int,
	machineCode string,
	machineName string,
) (*User, *UserToken, OperationResult) {
	// 查询用户是否存在或者被封禁
	var queryUsers []User
	queryResult := database.Where(&User{
		ID: userID,
	}).Find(&queryUsers)
	if queryResult.Error != nil {
		return nil, nil, OperationResult{
			Code:  errcode.DatabaseExecuteError,
			Error: queryResult.Error,
		}
	}
	if len(queryUsers) <= 0 {
		return nil, nil, OperationResult{
			Code:    errcode.ParamsError,
			Message: "User not exists",
		}
	}
	queryUser := queryUsers[0]
	if queryUser.Status == valueUserStatusForbidden {
		return nil, nil, OperationResult{
			Code:    errcode.UserForbidden,
			Message: "User is forbidden",
		}
	}
	// 查询 token 是否存在
	var queryUserTokens []UserToken
	queryResult = database.Where(&UserToken{
		UserID:      userID,
		Token:       refreshToken,
		Platform:    platform,
		MachineCode: machineCode,
		MachineName: machineName,
	})
	if queryResult.Error != nil {
		return nil, nil, OperationResult{
			Code:  errcode.DatabaseExecuteError,
			Error: queryResult.Error,
		}
	}
	if len(queryUserTokens) <= 0 {
		return nil, nil, OperationResult{
			Code:    errcode.AuthTokenInvalid,
			Message: "Token not exists",
		}
	}
	queryUserToken := queryUserTokens[0]
	// 判断 Token 是否过期
	currentTime := timeutils.Timestamp()
	if currentTime > queryUserToken.RefreshTokenExpireTime {
		// Refresh Token 过期，则删除信息并返回错误
		queryResult = database.Delete(&queryUserToken)
		if queryResult.Error != nil {
			log.Warming(tag, formatstring.String("Fail to remove expired refresh token: token=%s id=%d", queryUserToken.RefreshToken, queryUserToken.ID), zap.Error(queryResult.Error))
		}
		return nil, nil, OperationResult{
			Code:    errcode.AuthTokenExpired,
			Message: "Refresh Token Expired",
		}
	}
	// 生成新的 Token
	queryUserToken.Token = formatstring.GenerateToken()
	queryUserToken.TokenExpireTime = currentTime + apiconstuser.TimeTokenDuration
	queryUserToken.LastRefreshTime = currentTime
	queryResult = database.Save(&queryUserToken)
	if queryResult.Error != nil {
		log.Warming(tag, formatstring.String("Failed to update token: token=%s userID=%d", queryUserToken.Token, userID), zap.Error(queryResult.Error))
		return nil, nil, OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "Failed to update token",
			Error:   queryResult.Error,
		}
	}

	return &queryUser, &queryUserToken, OperationResult{Code: errcode.OK}
}
