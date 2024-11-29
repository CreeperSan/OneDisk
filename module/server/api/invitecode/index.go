package apiinvitecode

import (
	apiv1invitecode "OneDisk/module/server/api/invitecode/v1"
	"github.com/gin-gonic/gin"
)

func Register(server *gin.Engine) {
	apiv1invitecode.RegisterUserInviteCode(server)
}
