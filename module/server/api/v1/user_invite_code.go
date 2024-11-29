package apiv1user

import "github.com/gin-gonic/gin"

func registerUserInviteCode(server *gin.Engine) {

	/* 邀请码 - 创建 */
	server.POST("/api/user/v1/invite_code/generate", func(context *gin.Context) {

	})

	/* 邀请码 -  删除 */
	server.POST("/api/user/v1/invite_code/delete", func(context *gin.Context) {

	})

}
