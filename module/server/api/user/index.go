package apiuser

import (
	apiuserv1 "OneDisk/module/server/api/user/v1"
	"github.com/gin-gonic/gin"
)

func Register(server *gin.Engine) {
	apiuserv1.RegisterUserAuth(server)
	apiuserv1.RegisterUserLogin(server)
	apiuserv1.RegisterUserRegister(server)
}
