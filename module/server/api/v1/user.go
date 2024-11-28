package apiv1user

import (
	httpcode "OneDisk/definition/http_code"
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
		},
	)

	/* 认证 - 校验Token */
	server.POST(
		"/api/user/v1/auth/token",
		apimiddleware.Auth(),
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
