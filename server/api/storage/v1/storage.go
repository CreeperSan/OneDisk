package apistoragev1

import (
	defheader "OneDisk/def/header"
	defhttpcode "OneDisk/def/http_code"
	defstorage "OneDisk/def/storage"
	"OneDisk/lib/format/formatstring"
	"OneDisk/lib/log"
	"OneDisk/module/database"
	apimiddleware "OneDisk/server/api/middleware"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const tag = "ApiStorageV1"

func RegisterStorage(server *gin.Engine) {

	/* 添加存储 */
	server.POST(
		"/api/storage/v1/create",
		apimiddleware.AuthRequireAdminister(),
		func(context *gin.Context) {
			// 1、读取 Header 信息
			contextHeader, _ := context.Get(apimiddleware.KeyHeader)
			requestHeader, isInstance := contextHeader.(defheader.Header)
			if !isInstance {
				context.JSON(defhttpcode.InternalError, gin.H{
					"code": defhttpcode.InternalError,
					"msg":  "服务器内部错误，请稍后重试",
				})
				return
			}
			// 2、读取请求参数
			type RequestParams struct {
				Name   string `json:"name"`
				Avatar string `json:"avatar"`
				Type   int    `json:"type"`
			}
			var request RequestParams
			if err := context.BindJSON(&request); err != nil {
				context.JSON(defhttpcode.ParamsError, gin.H{
					"code": defhttpcode.ParamsError,
					"msg":  "参数错误",
				})
				return
			}
			// 3、校验参数
			if request.Type != database.ValueStorageTypePath {
				context.JSON(defhttpcode.ParamsError, gin.H{
					"code": defhttpcode.ParamsError,
					"msg":  "不支持的存储类型",
				})
			}
			// 4、存储参数配置读取并保存
			if request.Type == database.ValueStorageTypePath {
				// 4.1.1、本地存储策略
				type RequestStorageConfigLocalPath struct {
					Path string `json:"path"`
				}
				var requestStorageConfig RequestStorageConfigLocalPath
				if err := context.BindJSON(&requestStorageConfig); err != nil {
					context.JSON(defhttpcode.ParamsError, gin.H{
						"code": defhttpcode.ParamsError,
						"msg":  "参数错误",
					})
					return
				}
				if len(requestStorageConfig.Path) <= 0 {
					context.JSON(defhttpcode.ParamsError, gin.H{
						"code": defhttpcode.ParamsError,
						"msg":  "路径不能为空",
					})
					return
				}
				// 4.1.2、转换配置信息为 JSON 字符串
				configLocalPath := defstorage.ConfigLocalPath{
					Path: requestStorageConfig.Path,
				}
				configLocalPathJsonStr, err := json.Marshal(configLocalPath)
				if err != nil {
					context.JSON(defhttpcode.InternalError, gin.H{
						"code": defhttpcode.InternalError,
						"msg":  "服务器内部错误，请稍后重试",
					})
					return
				}
				// 4.1.3、保存存储策略配置
				insertStorage, result := database.StorageCreateAndSaveForLocalPath(
					request.Name,
					request.Avatar,
					requestHeader.UserID,
					request.Type,
					string(configLocalPathJsonStr),
				)
				if result.Code != defhttpcode.OK {
					context.JSON(result.Code, gin.H{
						"code": result.Code,
						"msg":  "操作失败",
					})
					return
				}
				// 4.1.4、默认建立创建人和与存储策略的对应关系
				_, result = database.StorageUserRelationCreateIfEmpty(requestHeader.UserID, insertStorage.ID)
				if result.Code != defhttpcode.OK {
					context.JSON(result.Code, gin.H{
						"code": result.Code,
						"msg":  "存储策略创建成功，但建立自身用户关联失败",
					})
					return
				}
				// 4.1.4、返回结果
				context.JSON(defhttpcode.OK, gin.H{
					"code": defhttpcode.OK,
					"msg":  "操作成功",
					"data": gin.H{
						"id":             insertStorage.ID,
						"create_user_id": insertStorage.CreateUserID,
						"name":           insertStorage.Name,
						"avatar":         insertStorage.Avatar,
						"type":           insertStorage.Type,
						"create_time":    insertStorage.CreateTime,
						"update_time":    insertStorage.UpdateTime,
						"path":           requestStorageConfig.Path,
						"config":         insertStorage.Config,
					},
				})
				return
			}
			context.JSON(defhttpcode.ParamsError, gin.H{
				"code": defhttpcode.ParamsError,
				"msg":  "不支持的存储类型",
			})
		},
	)

	/* 获取存储配置列表 */
	server.GET(
		"/api/storage/v1/list",
		apimiddleware.AuthToken(),
		func(context *gin.Context) {
			// 1、读取 Header 信息
			contextHeader, _ := context.Get(apimiddleware.KeyHeader)
			requestHeader, isInstance := contextHeader.(defheader.Header)
			if !isInstance {
				context.JSON(defhttpcode.InternalError, gin.H{
					"code": defhttpcode.InternalError,
					"msg":  "服务器内部错误，请稍后重试",
				})
				return
			}
			// 2、读取配置
			queryStorages, result := database.StorageUserRelationList(requestHeader.UserID)
			if result.Code != defhttpcode.OK {
				context.JSON(result.Code, gin.H{
					"code": result.Code,
					"msg":  "操作失败",
				})
				return
			}
			// 3、返回结果
			var storageList []gin.H
			for _, storage := range queryStorages {
				var configData map[string]interface{}
				errJson := json.Unmarshal([]byte(storage.Config), &configData)
				if errJson != nil {
					log.Warming(tag, formatstring.String("Unmarshal storage config failed, storageID=%d, configJson=%s", storage.ID, storage.Config), zap.Error(errJson))
					continue
				}
				storageList = append(storageList, gin.H{
					"id":             storage.ID,
					"create_user_id": storage.CreateUserID,
					"name":           storage.Name,
					"avatar":         storage.Avatar,
					"type":           storage.Type,
					"create_time":    storage.CreateTime,
					"update_time":    storage.UpdateTime,
					"config":         gin.H(configData),
				})
			}
			context.JSON(defhttpcode.OK, gin.H{
				"code": defhttpcode.OK,
				"msg":  "操作成功",
				"data": storageList,
			})
		},
	)

	/* 获取自己创建的存储配置列表 */
	server.GET(
		"/api/storage/v1/create/list",
		func(context *gin.Context) {
			// 1、读取 Header 信息
			contextHeader, _ := context.Get(apimiddleware.KeyHeader)
			requestHeader, isInstance := contextHeader.(defheader.Header)
			if !isInstance {
				context.JSON(defhttpcode.InternalError, gin.H{
					"code": defhttpcode.InternalError,
					"msg":  "服务器内部错误，请稍后重试",
				})
				return
			}
			// 2、读取配置
			queryStorages, result := database.StorageListByCreatUserID(requestHeader.UserID)
			if result.Code != defhttpcode.OK {
				context.JSON(result.Code, gin.H{
					"code": result.Code,
					"msg":  "操作失败",
				})
				return
			}
			// 3、返回结果
			var storageList []gin.H
			for _, storage := range queryStorages {
				var configData map[string]interface{}
				errJson := json.Unmarshal([]byte(storage.Config), &configData)
				if errJson != nil {
					log.Warming(tag, formatstring.String("Unmarshal storage config failed, storageID=%d, configJson=%s", storage.ID, storage.Config), zap.Error(errJson))
					continue
				}
				storageList = append(storageList, gin.H{
					"id":             storage.ID,
					"create_user_id": storage.CreateUserID,
					"name":           storage.Name,
					"avatar":         storage.Avatar,
					"type":           storage.Type,
					"create_time":    storage.CreateTime,
					"update_time":    storage.UpdateTime,
					"config":         gin.H(configData),
				})
			}
			context.JSON(defhttpcode.OK, gin.H{
				"code": defhttpcode.OK,
				"msg":  "操作成功",
				"data": storageList,
			})
		},
	)

}