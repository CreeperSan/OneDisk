package apimiddleware

import (
	httpcode "OneDisk/definition/http_code"
	apiconstheader "OneDisk/module/server/api/const/header"
	apimodel "OneDisk/module/server/api/const/model"
	"github.com/gin-gonic/gin"
	"strconv"
)

func HeaderConvert() gin.HandlerFunc {
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

		// 写入数据
		context.Set(apiconstheader.Token, apimodel.Header{
			Token:       headerToken,
			UserID:      headerUserIDInt,
			MachineCode: headerMachineCode,
			MachineName: headerMachineName,
			Platform:    headerPlatformInt,
		})

		context.Next()
	}
}
