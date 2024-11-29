package apiv1invitecode

import (
	httpcode "OneDisk/definition/http_code"
	"OneDisk/module/database"
	apimodel "OneDisk/module/server/api/const/model"
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
		// 1、读取 UserID
		contextHeader, _ := context.Get(apimiddleware.KeyHeader)
		requestHeader, isInstance := contextHeader.(apimodel.Header)
		if !isInstance {
			context.JSON(httpcode.InternalError, gin.H{
				"code": httpcode.InternalError,
				"msg":  "服务器内部错误，请稍后重试",
			})
			return
		}
		// 2、创建并保存邀请码
		insertInviteCode, result := database.InviteCodeCreateAndSaveForRegister(requestHeader.UserID)
		if result.Code != httpcode.OK {
			context.JSON(httpcode.InternalError, gin.H{
				"code": result.Code,
				"msg":  "操作失败，请重试",
			})
			return
		}
		// 3、返回数据
		context.JSON(httpcode.OK, gin.H{
			"code": httpcode.OK,
			"msg":  "操作成功",
			"data": gin.H{
				"id":           insertInviteCode.ID,
				"inviteCode":   insertInviteCode.Code,
				"expired_time": insertInviteCode.ExpiredTime,
			},
		})
	})

	/* 邀请码 -  删除 */
	requestGroup.POST("/delete", func(context *gin.Context) {

	})

}
