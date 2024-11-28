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
			context.JSON(httpcode.ParamsError, gin.H{
				"code": httpcode.ParamsError,
				"msg":  "操作失败，请重试",
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
