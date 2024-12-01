package apiuserv1

import (
	errcode "OneDisk/definition/err_code"
	httpcode "OneDisk/definition/http_code"
	"OneDisk/module/database"
	"OneDisk/server/api/const/model"
	apimiddleware2 "OneDisk/server/api/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterUserAuth(server *gin.Engine) {
	/* 认证 - 刷新 Token */
	server.POST(
		"/api/user/v1/auth/refresh",
		func(context *gin.Context) {
			// 使用新的 Token 替换旧的 Token
			// 1、读取 Header
			contextHeader, _ := context.Get(apimiddleware2.KeyHeader)
			requestHeader, isInstance := contextHeader.(apimodel.Header)
			if !isInstance {
				context.JSON(httpcode.InternalError, gin.H{
					"code": httpcode.InternalError,
					"msg":  "服务器内部错误，请稍后重试",
				})
				return
			}
			// 2、读取参数中的 refreshToken
			type RequestRefreshToken struct {
				RefreshToken string `json:"refreshToken"`
			}
			var request RequestRefreshToken
			if err := context.BindJSON(&request); err != nil {
				context.JSON(httpcode.ParamsError, gin.H{
					"code": httpcode.ParamsError,
					"msg":  "操作失败，请重试",
				})
				return
			}
			// 3、更新到数据库
			_, userToken, result := database.UserTokenRefresh(
				requestHeader.UserID,
				request.RefreshToken,
				requestHeader.Platform,
				requestHeader.MachineCode,
				requestHeader.MachineName,
			)
			if result.Code != errcode.OK {
				context.JSON(httpcode.ParamsError, gin.H{
					"code": result.Code,
					"msg":  "操作失败，请重试",
				})
				return
			}
			// 4、返回新的 Token
			context.JSON(httpcode.OK, gin.H{
				"code": httpcode.OK,
				"msg":  "操作成功",
				"data": gin.H{
					"token":        userToken.Token,
					"refreshToken": userToken.RefreshToken,
				},
			})
		},
	)

	/* 认证 - 校验Token */
	server.POST(
		"/api/user/v1/auth/token",
		apimiddleware2.AuthToken(),
		func(context *gin.Context) {
			// 中间件已经处理，可以直接返回
			context.JSON(httpcode.OK, gin.H{
				"msg": "操作成功",
			})
		},
	)
}
