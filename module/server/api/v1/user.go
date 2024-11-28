package apiv1user

import (
	errcode "OneDisk/definition/err_code"
	httpcode "OneDisk/definition/http_code"
	"OneDisk/lib/format/formatstring"
	timeutils "OneDisk/lib/utils/time"
	"OneDisk/module/database"
	apimodel "OneDisk/module/server/api/const/model"
	apiconstuser "OneDisk/module/server/api/const/user"
	apimiddleware "OneDisk/module/server/api/middleware"
	"github.com/gin-gonic/gin"
)

func Register(server *gin.Engine) {

	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////// 用户认证

	/* 认证 - 刷新 Token */
	server.POST(
		"/api/user/v1/auth/refresh",
		func(context *gin.Context) {
			// 使用新的 Token 替换旧的 Token
			// 1、读取 Header
			contextHeader, _ := context.Get(apimiddleware.KeyHeader)
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
		apimiddleware.AuthToken(),
		func(context *gin.Context) {
			// 中间件已经处理，可以直接返回
			context.JSON(httpcode.OK, gin.H{
				"msg": "操作成功",
			})
		},
	)

	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////// 用户登录

	/* 登录 */
	server.POST("/api/user/v1/login", func(context *gin.Context) {
		// 检查 Header
		contextHeader, _ := context.Get(apimiddleware.KeyHeader)
		requestHeader, isInstance := contextHeader.(apimodel.Header)
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

	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////// 用户注册

	/* 注册 */
	server.POST("/api/user/v1/register", func(context *gin.Context) {

	})

	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////// 用户邀请码

	/* 邀请码 - 创建 */
	server.POST("/api/user/v1/invite_code/generate", func(context *gin.Context) {

	})

	/* 邀请码 -  删除 */
	server.POST("/api/user/v1/invite_code/delete", func(context *gin.Context) {

	})

}
