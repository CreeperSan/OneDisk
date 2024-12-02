package apimiddleware

import (
	"OneDisk/def/header"
	httpcode "OneDisk/def/http_code"
	"github.com/gin-gonic/gin"
	"strconv"
)

func HeaderConvert() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 读取 Header
		headerToken := context.GetHeader(defheader.Token)
		headerMachineCode := context.GetHeader(defheader.MachineCode)
		headerMachineName := context.GetHeader(defheader.MachineName)
		headerPlatform := context.GetHeader(defheader.Platform)
		headerUserID := context.GetHeader(defheader.UserID)

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
		context.Set(defheader.Token, defheader.Header{
			Token:       headerToken,
			UserID:      headerUserIDInt,
			MachineCode: headerMachineCode,
			MachineName: headerMachineName,
			Platform:    headerPlatformInt,
		})

		context.Next()
	}
}
