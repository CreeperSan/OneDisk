package apiv1invitecode

import (
	apimiddleware "OneDisk/module/server/api/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterUserInviteCode(server *gin.Engine) {
	requestGroup := server.Group("/api/user/v1/invite_code")

	// 邀请码都需要已登录用户
	requestGroup.Use(apimiddleware.AuthToken())
	// 邀请码都需要管理员操作
	requestGroup.Use(apimiddleware.AuthRequireAdminister())

	/* 邀请码 - 创建 */
	requestGroup.POST("/generate", func(context *gin.Context) {

	})

	/* 邀请码 -  删除 */
	requestGroup.POST("/delete", func(context *gin.Context) {

	})

}
