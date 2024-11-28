package errcode

const (
	OK = 0

	DatabaseExecuteError = 1
	ParamsError          = 2

	UserNotExist  = 10001
	UserForbidden = 10002

	AuthTokenInvalid = 20001
	AuthTokenExpired = 20002
)
