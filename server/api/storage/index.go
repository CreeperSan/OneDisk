package apistorage

import (
	apistoragev1 "OneDisk/server/api/storage/v1"
	"github.com/gin-gonic/gin"
)

func Register(server *gin.Engine) {
	apistoragev1.RegisterStorage(server)
}
