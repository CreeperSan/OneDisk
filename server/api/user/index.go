package apiuser

import (
	apiuserv2 "OneDisk/server/api/user/v1"
	"github.com/gin-gonic/gin"
)

func Register(server *gin.Engine) {
	apiuserv2.RegisterUserAuth(server)
	apiuserv2.RegisterUserLogin(server)
	apiuserv2.RegisterUserRegister(server)
}
