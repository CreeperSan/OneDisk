package apiuserv1

import (
	errcode "OneDisk/def/err_code"
	defheader "OneDisk/def/header"
	httpcode "OneDisk/def/http_code"
	"OneDisk/module/database"
	"OneDisk/server/api/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterUserLogin(server *gin.Engine) {
	/* 登录 */
	server.POST("/api/user/v1/login", func(context *gin.Context) {
		// 检查 Header
		contextHeader, _ := context.Get(apimiddleware.KeyHeader)
		requestHeader, isInstance := contextHeader.(defheader.Header)
		if !isInstance {
			context.JSON(httpcode.InternalError, gin.H{
				"code": httpcode.InternalError,
				"msg":  "服务器内部错误，请稍后重试",
			})
			return
		}
		// 检查请求参数
		type RequestLogin struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		var request RequestLogin
		if err := context.BindJSON(&request); err != nil {
			context.JSON(httpcode.ParamsError, gin.H{
				"code": httpcode.ParamsError,
				"msg":  "操作失败，请重试",
			})
			return
		}
		// 尝试用户名密码登录
		resultUser, result := database.UserValidationByUsername(request.Username, request.Password)
		if result.Code == errcode.DatabaseExecuteError {
			context.JSON(httpcode.InternalError, gin.H{
				"code": result.Code,
				"msg":  "服务器内部错误，请重试",
			})
			return
		}
		if resultUser == nil {
			context.JSON(httpcode.ParamsError, gin.H{
				"code": result.Code,
				"msg":  "用户名或密码错误,请重试",
			})
			return
		}
		// 防止重复登录，先删除旧的 token 记录
		result = database.UserTokenRemove(
			requestHeader.Platform,
			requestHeader.MachineCode,
			requestHeader.MachineName,
		)
		if result.Code != errcode.OK {
			context.JSON(httpcode.InternalError, gin.H{
				"code": result.Code,
				"msg":  "服务器内部错误，请重试",
			})
			return
		}
		// 登录成功，生成并保存 Token
		resultUser, resultUserToken, result := database.UserTokenCreateAndInsert(
			resultUser.ID,
			requestHeader.Platform,
			requestHeader.MachineCode,
			requestHeader.MachineName,
		)
		if result.Code != errcode.OK {
			context.JSON(httpcode.InternalError, gin.H{
				"code": result.Code,
				"msg":  "服务器内部错误，请重试",
			})
			return
		}
		// 返回数据
		context.JSON(httpcode.OK, gin.H{
			"code": httpcode.OK,
			"msg":  "登录成功",
			"data": gin.H{
				"userID":       resultUser.ID,
				"username":     resultUser.Username,
				"nickname":     resultUser.Nickname,
				"avatar":       resultUser.Avatar,
				"email":        resultUser.Email,
				"phone":        resultUser.Phone,
				"type":         resultUser.Type,
				"status":       resultUser.Status,
				"token":        resultUserToken.Token,
				"refreshToken": resultUserToken.RefreshToken,
			},
		})
	})

}
