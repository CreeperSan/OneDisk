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
	"path/filepath"
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

	/* 文件 - 获取文件列表 */
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

	/* 文件 - 创建文件夹 */
	server.PUT(
		"/api/file/v1/create_directory",
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
				return
			}
			// 2、读取参数
			type RequestData struct {
				Path string `json:"path"`
				Name string `json:"name"`
			}
			var requestData RequestData
			err := context.BindJSON(&requestData)
			if err != nil {
				context.JSON(httpcode.ParamsError, gin.H{
					"code": httpcode.ParamsError,
					"msg":  "请求参数错误",
				})
				return
			}
			// 3、调用平台接口
			createDirectory, result := storage.CreateDirectory(requestData.Path + "/" + requestData.Name)
			if result.Code != errcode.OK {
				context.JSON(httpcode.ParamsError, gin.H{
					"code": result.Code,
					"msg":  "服务器内部错误，请稍后重试",
				})
				return
			}
			if createDirectory == nil {
				context.JSON(httpcode.ParamsError, gin.H{
					"code": httpcode.ParamsError,
					"msg":  "创建失败",
				})
				return
			}
			// 4. 返回结果
			data, err := convertFileToJson(*createDirectory)
			if err != nil {
				context.JSON(httpcode.ParamsError, gin.H{
					"code": httpcode.ParamsError,
					"msg":  "创建失败",
				})
				return
			}
			context.JSON(httpcode.OK, gin.H{
				"code": httpcode.OK,
				"msg":  "操作成功",
				"data": data,
			})
		},
	)

	/* 文件 - 创建文件 */
	server.PUT(
		"/api/file/v1/create_file",
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
				return
			}
			// 2、读取参数
			type RequestData struct {
				Path string `json:"path"`
				Name string `json:"name"`
			}
			var requestData RequestData
			err := context.BindJSON(&requestData)
			if err != nil {
				context.JSON(httpcode.ParamsError, gin.H{
					"code": httpcode.ParamsError,
					"msg":  "请求参数错误",
				})
				return
			}
			// 3、调用平台接口
			createDirectory, result := storage.CreateFile(requestData.Path + "/" + requestData.Name)
			if result.Code != errcode.OK {
				context.JSON(httpcode.ParamsError, gin.H{
					"code": result.Code,
					"msg":  "服务器内部错误，请稍后重试",
				})
				return
			}
			if createDirectory == nil {
				context.JSON(httpcode.ParamsError, gin.H{
					"code": httpcode.ParamsError,
					"msg":  "创建失败",
				})
				return
			}
			// 4. 返回结果
			data, err := convertFileToJson(*createDirectory)
			if err != nil {
				context.JSON(httpcode.ParamsError, gin.H{
					"code": httpcode.ParamsError,
					"msg":  "创建失败",
				})
				return
			}
			context.JSON(httpcode.OK, gin.H{
				"code": httpcode.OK,
				"msg":  "操作成功",
				"data": data,
			})
		},
	)

	/* 文件 - 删除文件 */
	server.DELETE(
		"/api/file/v1/delete",
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
				return
			}
			// 2、读取参数
			type RequestData struct {
				Path string `json:"path"`
				Name string `json:"name"`
			}
			var requestData RequestData
			err := context.BindJSON(&requestData)
			if err != nil {
				context.JSON(httpcode.ParamsError, gin.H{
					"code": httpcode.ParamsError,
					"msg":  "请求参数错误",
				})
				return
			}
			// 3、调用平台接口
			result := storage.Delete(requestData.Path + "/" + requestData.Name)
			if result.Code != errcode.OK {
				context.JSON(httpcode.ParamsError, gin.H{
					"code": result.Code,
					"msg":  "服务器内部错误，请稍后重试",
				})
				return
			}
			context.JSON(httpcode.OK, gin.H{
				"code": httpcode.OK,
				"msg":  "操作成功",
			})
		},
	)

	/* 文件 - 移动文件 */
	server.POST(
		"/api/file/v1/move",
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
				return
			}
			// 2、读取参数
			type RequestData struct {
				FromPath string `json:"from_path"`
				ToPath   string `json:"to_path"`
			}
			var requestData RequestData
			err := context.BindJSON(&requestData)
			if err != nil {
				context.JSON(httpcode.ParamsError, gin.H{
					"code": httpcode.ParamsError,
					"msg":  "请求参数错误",
				})
				return
			}
			// 3、调用平台接口
			createDirectory, result := storage.Move(requestData.FromPath, requestData.ToPath)
			if result.Code != errcode.OK {
				context.JSON(httpcode.ParamsError, gin.H{
					"code": result.Code,
					"msg":  "服务器内部错误，请稍后重试",
				})
				return
			}
			if createDirectory == nil {
				context.JSON(httpcode.ParamsError, gin.H{
					"code": httpcode.ParamsError,
					"msg":  "创建失败",
				})
				return
			}
			// 4. 返回结果
			data, err := convertFileToJson(*createDirectory)
			if err != nil {
				context.JSON(httpcode.ParamsError, gin.H{
					"code": httpcode.ParamsError,
					"msg":  "创建失败",
				})
				return
			}
			context.JSON(httpcode.OK, gin.H{
				"code": httpcode.OK,
				"msg":  "操作成功",
				"data": data,
			})
		},
	)

	/* 文件 - 上传文件 */
	server.POST(
		"/api/file/v1/upload",
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
				return
			}
			// 2、读取文件
			requestFile, err := context.FormFile("file")
			if err != nil {
				context.JSON(httpcode.ParamsError, gin.H{
					"code": httpcode.ParamsError,
					"msg":  "请求参数错误",
				})
				return
			}
			// 3、读取文件路径
			type RequestParams struct {
				Path string `json:"path"`
			}
			var requestParams RequestParams
			err = context.BindJSON(&requestParams)
			if err != nil {
				context.JSON(httpcode.ParamsError, gin.H{
					"code": httpcode.ParamsError,
					"msg":  "请求参数错误",
				})
				return
			}
			requestPathAbsolute, err := filepath.Abs(requestParams.Path)
			if err != nil {
				context.JSON(httpcode.ParamsError, gin.H{
					"code": httpcode.ParamsError,
					"msg":  "请求参数错误",
				})
				return
			}
			// 4、调用平台接口
			resultFile, result := storage.Upload(context, requestFile, requestPathAbsolute)
			if result.Code != errcode.OK {
				context.JSON(httpcode.ParamsError, gin.H{
					"code": result.Code,
					"msg":  "上传文件失败，请稍后重试",
				})
				return
			}
			// 5、返回上传的文件信息
			data, err := convertFileToJson(*resultFile)
			if err != nil {
				context.JSON(httpcode.ParamsError, gin.H{
					"code": httpcode.ParamsError,
					"msg":  "上传文件失败，请稍后重试",
				})
				return
			}
			context.JSON(httpcode.OK, gin.H{
				"code": httpcode.OK,
				"msg":  "操作成功",
				"data": data,
			})
		},
	)

	/* 文件 - 生成下载链接 */
	server.POST(
		"/api/file/v1/generate_download_url",
		apifilemiddleware.StorageGetPlatformInterface(),
		func(context *gin.Context) {

		},
	)

	/* 文件 - 透过下载链接下载文件 */
	server.POST(
		"/api/file/v1/download",
		apifilemiddleware.StorageGetPlatformInterface(),
		func(context *gin.Context) {

		},
	)

}
