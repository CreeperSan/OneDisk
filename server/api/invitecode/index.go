package apiinvitecode

import (
	"OneDisk/server/api/invitecode/v1"
	"github.com/gin-gonic/gin"
)

func Register(server *gin.Engine) {
	apiv1invitecode.RegisterUserInviteCode(server)
}
