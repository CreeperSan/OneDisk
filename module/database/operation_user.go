package database

import (
	errcode "OneDisk/definition/err_code"
	"OneDisk/lib/format/formatstring"
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
