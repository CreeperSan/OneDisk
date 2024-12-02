package apifile

import (
	apifilev1 "OneDisk/server/api/file/v1"
	"github.com/gin-gonic/gin"
)

func Register(server *gin.Engine) {
	apifilev1.RegisterFile(server)
}
