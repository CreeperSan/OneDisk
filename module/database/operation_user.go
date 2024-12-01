package database

import (
	errcode "OneDisk/definition/err_code"
	"OneDisk/lib/format/formatstring"
	"time"
)

// UserFindUser
// 查询用户
func UserFindUser(userID int64) (*User, OperationResult) {
	if userID <= 0 {
		return nil, OperationResult{
			Code:    errcode.ParamsError,
			Message: "UserID not valid",
		}
	}

	var queryUsers []User
	queryResult := database.Where(&User{ID: userID}).Find(&queryUsers)
	if queryResult.Error != nil {
		return nil, OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "Error occurred while querying user in UserFindUser()",
			Error:   queryResult.Error,
		}
	}
	if len(queryUsers) <= 0 {
		return nil, OperationResult{
			Code:    errcode.UserNotExist,
			Message: "User not exist",
		}
	}

	return &queryUsers[0], OperationResult{Code: errcode.OK}
}

// UserFindUserByUsername
// 通过用户名查询用户
func UserFindUserByUsername(username string) (*User, OperationResult) {
	if len(username) <= 0 {
		return nil, OperationResult{
			Code:    errcode.ParamsError,
			Message: "Username can not be empty",
		}
	}

	var queryUsers []User
	queryResult := database.Where(&User{Username: username}).Find(&queryUsers)
	if queryResult.Error != nil {
		return nil, OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "Error occurred while querying user in UserFindUserByUsername()",
			Error:   queryResult.Error,
		}
	}
	if len(queryUsers) <= 0 {
		return nil, OperationResult{
			Code:    errcode.UserNotExist,
			Message: "User not exist",
		}
	}

	return &queryUsers[0], OperationResult{Code: errcode.OK}
}

func UserValidationByUsername(username string, password string) (*User, OperationResult) {
	if len(username) <= 0 || len(password) <= 0 {
		return nil, OperationResult{
			Code:    errcode.ParamsError,
			Message: "Username and password can not be empty",
		}
	}

	passwordEncode := formatstring.Password(password)

	// 查询用户
	var queryUsers []User
	queryResult := database.Where(&User{
		Username: username,
		Password: passwordEncode,
	}).Find(&queryUsers)
	if queryResult.Error != nil {
		return nil, OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "Error occurred while querying user in UserValidationByUsername()",
			Error:   queryResult.Error,
		}
	}
	if len(queryUsers) <= 0 {
		return nil, OperationResult{
			Code:    errcode.UserNotExist,
			Message: "User not exist",
		}
	}

	return &queryUsers[0], OperationResult{
		Code: errcode.OK,
	}
}

func UserCreateAndSave(username string, password string, usertype int) (*User, OperationResult) {
	// 1、参数检查
	if len(username) <= 0 || len(password) <= 0 {
		return nil, OperationResult{
			Code:    errcode.ParamsError,
			Message: "Username and password can not be empty",
		}
	}
	// 2、检查用户是否已经存在
	queryUser, result := UserFindUserByUsername(username)
	if result.Code != errcode.OK {
		return nil, result
	}
	if queryUser != nil {
		return nil, OperationResult{
			Code:    errcode.UserExist,
			Message: "User already",
		}
	}
	// 3、创建用户
	currentTimestamp := timeutils.Timestamp()
	passwordEncode := formatstring.Password(password)
	newUser := User{
		Username:   username,
		Password:   passwordEncode,
		Email:      "",
		Nickname:   username,
		Avatar:     "",
		Phone:      "",
		CreateTime: currentTimestamp,
		Type:       usertype,
		Status:     ValueUserStatusActive,
	}
	// 4、保存用户
	saveResult := database.Create(&newUser)
	if saveResult.Error != nil {
		return nil, OperationResult{
			Code:    errcode.DatabaseExecuteError,
			Message: "Error occurred while saving user in UserCreateAndSave()",
			Error:   saveResult.Error,
		}
	}
	return &newUser, OperationResult{Code: errcode.OK}
}
