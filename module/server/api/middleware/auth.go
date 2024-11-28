package apimiddleware

import (
	errcode "OneDisk/definition/err_code"
	httpcode "OneDisk/definition/http_code"
	"OneDisk/module/database"
	apiconstheader "OneDisk/module/server/api/const/header"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 读取 Header
		headerToken := context.GetHeader(apiconstheader.Token)
		headerMachineCode := context.GetHeader(apiconstheader.MachineCode)
		headerMachineName := context.GetHeader(apiconstheader.MachineName)
		headerPlatform := context.GetHeader(apiconstheader.Platform)
		headerUserID := context.GetHeader(apiconstheader.UserID)
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

		validationUser, validationToken, validationResult := database.UserTokenValidation(headerUserIDInt, headerToken, headerMachineCode, headerMachineName, headerPlatformInt)

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
