package apimiddleware

import (
	errcode "OneDisk/definition/err_code"
	httpcode "OneDisk/definition/http_code"
	"OneDisk/module/database"
	apimodel "OneDisk/module/server/api/const/model"
	"github.com/gin-gonic/gin"
)

func AuthToken() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 读取 Header
		contextHeader, _ := context.Get(KeyHeader)
		requestHeader, isInstance := contextHeader.(apimodel.Header)
		if !isInstance {
			context.JSON(httpcode.InternalError, gin.H{
				"code": httpcode.InternalError,
				"msg":  "服务器内部错误，请重试",
			})
			context.Abort()
			return
		}

		if requestHeader.Token == "" || requestHeader.MachineCode == "" || requestHeader.MachineName == "" || requestHeader.Platform <= 0 || requestHeader.UserID <= 0 {
			context.JSON(httpcode.ParamsError, gin.H{
				"code": httpcode.ParamsError,
				"msg":  "操作失败，请重试",
			})
			context.Abort()
			return
		}

		validationUser, validationToken, validationResult := database.UserTokenValidation(requestHeader.UserID, requestHeader.Token, requestHeader.MachineCode, requestHeader.MachineName, requestHeader.Platform)

		if validationResult.Error != nil {
			context.JSON(httpcode.InternalError, gin.H{
				"code": validationResult.Code,
				"msg":  "服务器内部出错，请稍后重试",
			})
			context.Abort()
			return
		}

		if validationUser == nil || validationToken == nil || validationResult.Code != errcode.OK {
			context.JSON(httpcode.Unauthorized, gin.H{
				"code": httpcode.Unauthorized,
				"msg":  "登录信息过期，请重新登录",
			})
			context.Abort()
			return
		}

		// 储存信息
		context.Set(KeyUser, validationUser)
		context.Set(KeyUserToken, validationToken)

		// 继续后面处理
		context.Next()
	}
}

// AuthRequireAdminister
// 标记这个方法需要管理员权限才能调用
// 调用这个的前提依赖是需要先调用 AuthToken
func AuthRequireAdminister() gin.HandlerFunc {
	return func(context *gin.Context) {
		contextUser, _ := context.Get(KeyUser)
		requestUser, isInstance := contextUser.(*database.User)
		if !isInstance {
			context.JSON(httpcode.InternalError, gin.H{
				"code": httpcode.InternalError,
				"msg":  "服务器内部错误，请稍后重试",
			})
			return
		}
		if requestUser.Type != database.ValueUserTypeAdmin {
			context.JSON(httpcode.Forbidden, gin.H{
				"code": httpcode.Forbidden,
				"msg":  "您没有足够的权限操作",
			})
			context.Abort()
		}
	}
}
