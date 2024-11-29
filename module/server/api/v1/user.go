package apiv1user

import (
	"github.com/gin-gonic/gin"
)

func Register(server *gin.Engine) {
	registerUserAuth(server)
	registerUserLogin(server)
	registerUserRegister(server)
	registerUserInviteCode(server)
}
