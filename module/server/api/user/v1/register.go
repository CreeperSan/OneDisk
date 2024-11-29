package apiuserv1

import "github.com/gin-gonic/gin"

func RegisterUserRegister(server *gin.Engine) {

	/* 注册 */
	server.POST("/api/user/v1/register", func(context *gin.Context) {

	})

}
