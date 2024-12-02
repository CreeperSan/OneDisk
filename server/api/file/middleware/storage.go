package apifilemiddleware

import (
	errcode "OneDisk/def/err_code"
	defheader "OneDisk/def/header"
	httpcode "OneDisk/def/http_code"
	defstorage "OneDisk/def/storage"
	"OneDisk/module/database"
	apimiddleware "OneDisk/server/api/middleware"
	"OneDisk/storage"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
)

// StorageGetPlatformInterface
// 获取存储平台接口对象, 前提是已经经过了 Header 中间件
func StorageGetPlatformInterface() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 1、先读取 Header
		contextHeader, _ := context.Get(apimiddleware.KeyHeader)
		requestHeader, isInstance := contextHeader.(defheader.Header)
		if !isInstance {
			context.JSON(500, gin.H{
				"code": 500,
				"msg":  "服务器内部错误，请稍后重试",
			})
			context.Abort()
			return
		}
		// 2、先读取类型
		var requestStorageType int = database.ValueStorageTypeUndefined
		var requestStorageID int64 = 0
		if context.Request.Method == "GET" {
			queryStorageType := context.Query("type")
			queryStorageTypeInt, err := strconv.Atoi(queryStorageType)
			if err != nil {
				context.JSON(400, gin.H{
					"code": 400,
					"msg":  "参数错误",
				})
				context.Abort()
				return
			}
			queryStorageID := context.Query("storage_id")
			queryStorageIDInt, err := strconv.Atoi(queryStorageID)
			if err != nil {
				context.JSON(400, gin.H{
					"code": 400,
					"msg":  "参数错误",
				})
				context.Abort()
				return
			}
			requestStorageType = queryStorageTypeInt
			requestStorageID = int64(queryStorageIDInt)
		} else if context.Request.Method == "POST" {
			type RequestData struct {
				StorageID int `json:"storage_id"`
				Type      int `json:"type"`
			}
			var requestData RequestData
			err := context.BindJSON(&requestData)
			if err != nil {
				context.JSON(400, gin.H{
					"code": 400,
					"msg":  "参数错误",
				})
				context.Abort()
				return
			}
			requestStorageType = requestData.Type
			requestStorageID = int64(requestData.StorageID)
		} else {
			context.JSON(400, gin.H{
				"code": 400,
				"msg":  "暂未支持的请求方法",
			})
			context.Abort()
			return
		}
		// 3、再获取平台接口对象
		if requestStorageType == database.ValueStorageTypePath {
			// 3.1、读取配置
			queryStorageUserRelation, result := database.StorageUserRelationFindByUserIDAndStorageID(requestHeader.UserID, requestStorageID)
			if result.Code != errcode.OK || queryStorageUserRelation == nil {
				context.JSON(httpcode.ParamsError, gin.H{
					"code": httpcode.ParamsError,
					"msg":  "存储策略不存在",
				})
				context.Abort()
				return
			}
			queryStorage, result := database.StorageFind(requestStorageID)
			if result.Code != errcode.OK || queryStorage == nil {
				context.JSON(httpcode.ParamsError, gin.H{
					"code": httpcode.ParamsError,
					"msg":  "存储策略不存在",
				})
				context.Abort()
				return
			}
			var configLocalPath defstorage.ConfigLocalPath
			err := json.Unmarshal([]byte(queryStorage.Config), &configLocalPath)
			if err != nil {
				context.JSON(httpcode.InternalError, gin.H{
					"code": httpcode.InternalError,
					"msg":  "存储策略配置有误",
				})
				context.Abort()
				return
			}
			// 3.2、实例化
			prefabStorage := storage.PlatformInterfaceLocal{
				Root: configLocalPath.Path,
			}
			context.Set(apimiddleware.KeyStorage, prefabStorage)
			// 2.3、继续请求
			context.Next()
		} else {
			context.JSON(400, gin.H{
				"code": 400,
				"msg":  "暂未支持的存储策略",
			})
			context.Abort()
			return
		}
	}
}
