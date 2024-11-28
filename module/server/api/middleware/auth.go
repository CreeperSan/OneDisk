package apimiddleware

import (
	httpcode "OneDisk/definition/http_code"
	"OneDisk/module/database"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 读取 Header
		headerToken := context.GetHeader("one-disk-token")
		headerMachineCode := context.GetHeader("one-disk-machine-code")
		headerMachineName := context.GetHeader("one-disk-machine-name")
		headerPlatform := context.GetHeader("one-disk-platform")
		headerUserID := context.GetHeader("one-disk-user-id")
		// 参数转换
		headerPlatformInt, err := strconv.Atoi(headerPlatform)
		if err != nil {
			context.JSON(httpcode.ParamsError, gin.H{
				"code": httpcode.ParamsError,
				"msg":  "操作失败，请重试",
			})
			context.Abort()
			return
		}
		headerUserIDInt, err := strconv.ParseInt(headerUserID, 10, 64)
		if err != nil {
			context.JSON(httpcode.ParamsError, gin.H{
				"code": httpcode.ParamsError,
				"msg":  "操作失败，请重试",
			})
			context.Abort()
			return
		}

		if headerToken == "" || headerMachineCode == "" || headerMachineName == "" || headerPlatformInt <= 0 || headerUserIDInt <= 0 {
			context.JSON(httpcode.ParamsError, gin.H{
				"code": httpcode.ParamsError,
				"msg":  "操作失败，请重试",
			})
			context.Abort()
			return
		}

		validationUser, validationToken, validationError := database.UserTokenValidation(headerUserIDInt, headerToken, headerMachineCode, headerMachineName, headerPlatformInt)
		if validationError != nil {
			context.JSON(httpcode.ParamsError, gin.H{
				"code": validationError.Code,
				"msg":  validationError.Message,
			})
			context.Abort()
			return
		}

		if validationUser == nil || validationToken == nil {
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
