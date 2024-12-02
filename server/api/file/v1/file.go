package apifilev1

import (
	errcode "OneDisk/def/err_code"
	httpcode "OneDisk/def/http_code"
	"OneDisk/lib/log"
	apifilemiddleware "OneDisk/server/api/file/middleware"
	apimiddleware "OneDisk/server/api/middleware"
	storage2 "OneDisk/storage"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const tag = "ApiFileV1"

func convertFileToJson(file storage2.File) (gin.H, error) {
	data, err := json.Marshal(file)
	if err != nil {
		log.Warming(tag, "Error occur while convertFileToJson", zap.Error(err))
		return nil, err
	}
	var result gin.H
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Warming(tag, "Error occur while convertFileToJson", zap.Error(err))
		return nil, err
	}
	return result, nil
}

func RegisterFile(server *gin.Engine) {

	server.GET(
		"/api/file/v1/list",
		apifilemiddleware.StorageGetPlatformInterface(),
		func(context *gin.Context) {
			// 1、读取配置
			contextHeader, _ := context.Get(apimiddleware.KeyStorage)
			storage, isInstance := contextHeader.(storage2.PlatformInterface)
			if !isInstance {
				context.JSON(500, gin.H{
					"code": 500,
					"msg":  "服务器内部错误，请稍后重试",
				})
				context.Abort()
				return
			}
			// 2、读取参数
			requestPath := context.Query("path")
			if len(requestPath) <= 0 {
				requestPath = "/"
			}
			// 3、调用平台接口
			storageFileList, result := storage.List(requestPath)
			if result.Code != errcode.OK {
				context.JSON(httpcode.InternalError, gin.H{
					"code": result.Code,
					"msg":  "服务器内部错误，请稍后重试",
				})
				return
			}
			// 4、返回结果
			var dataFileList []gin.H
			for _, storageFile := range storageFileList {
				data, err := convertFileToJson(storageFile)
				if err != nil {
					continue
				}
				dataFileList = append(dataFileList, data)
			}
			context.JSON(httpcode.OK, gin.H{
				"code": httpcode.OK,
				"msg":  "操作成功",
				"data": dataFileList,
			})
		},
	)

}
