package apiuserv1

import (
	errcode "OneDisk/definition/err_code"
	httpcode "OneDisk/definition/http_code"
	"OneDisk/module/database"
	apiconstinvitecode "OneDisk/server/api/const/invitecode"
	apimiddleware "OneDisk/server/api/middleware"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func RegisterUserRegister(server *gin.Engine) {

	/* 注册 */
	server.POST("/api/user/v1/register", apimiddleware.HeaderConvert(), func(context *gin.Context) {
		// 1、获取并检查请求参数
		type RequestRegister struct {
			Username   string `json:"username"`
			Password   string `json:"password"`
			InviteCode string `json:"invite_code"`
		}
		var request RequestRegister
		if err := context.BindJSON(&request); err != nil {
			context.JSON(httpcode.ParamsError, gin.H{
				"code": httpcode.ParamsError,
				"msg":  "操作失败，请重试",
			})
			return
		}
		if len(request.Username) <= 0 || len(request.Password) <= 0 {
			context.JSON(httpcode.ParamsError, gin.H{
				"code": httpcode.ParamsError,
				"msg":  "用户名或密码不能为空",
			})
			return
		}
		if len(request.InviteCode) <= 0 {
			context.JSON(httpcode.ParamsError, gin.H{
				"code": httpcode.ParamsError,
				"msg":  "邀请码不能为空",
			})
			return
		}
		// 2、检查用户是否注册
		queryUser, result := database.UserFindUserByUsername(request.Username)
		if result.Code != errcode.OK {
			context.JSON(httpcode.InternalError, gin.H{
				"code": result.Code,
				"msg":  "服务器内部错误，请稍后重试",
			})
			return
		}
		if queryUser != nil {
			context.JSON(httpcode.ParamsError, gin.H{
				"code": httpcode.ParamsError,
				"msg":  "用户名已存在，请更换用户名称",
			})
			return
		}
		// 3、检查邀请码是否有效
		queryInviteCode, result := database.InviteCodeFindByCode(request.InviteCode)
		if result.Code != errcode.OK {
			context.JSON(httpcode.InternalError, gin.H{
				"code": result.Code,
				"msg":  "服务器内部错误，请稍后重试",
			})
			return
		}
		if queryInviteCode == nil {
			context.JSON(httpcode.ParamsError, gin.H{
				"code": httpcode.ParamsError,
				"msg":  "邀请码无效，请检查邀请码是否输入正确",
			})
			return
		}
		if queryInviteCode.Status != database.ValueInviteCodeStatusNotUse {
			if queryInviteCode.Status == database.ValueInviteCodeStatusInvalid {
				context.JSON(httpcode.ParamsError, gin.H{
					"code": httpcode.ParamsError,
					"msg":  "邀请码已失效，请更换邀请码",
				})
				return
			}
			if queryInviteCode.Status == database.ValueInviteCodeStatusUsed {
				context.JSON(httpcode.ParamsError, gin.H{
					"code": httpcode.ParamsError,
					"msg":  "邀请码已使用，请更换邀请码",
				})
				return
			}
			context.JSON(httpcode.ParamsError, gin.H{
				"code": httpcode.ParamsError,
				"msg":  "邀请码已不可用，请更换邀请码",
			})
			return
		}
		// 4、检查邀请码的用户是否封禁
		queryUser, result = database.UserFindUser(queryInviteCode.FromUserID)
		if result.Code != errcode.OK {
			context.JSON(httpcode.InternalError, gin.H{
				"code": result.Code,
				"msg":  "服务器内部错误，请稍后重试",
			})
			return
		}
		if queryUser == nil {
			context.JSON(httpcode.ParamsError, gin.H{
				"code": httpcode.ParamsError,
				"msg":  "邀请码已失效，请更换邀请码",
			})
			return
		}
		if queryUser.Status != database.ValueUserStatusActive {
			context.JSON(httpcode.ParamsError, gin.H{
				"code": httpcode.ParamsError,
				"msg":  "邀请码已失效，请更换邀请码",
			})
		}
		// 5、解析邀请码的额外信息
		var extraRegister apiconstinvitecode.ExtraForRegister
		if len(queryInviteCode.Extra) > 0 {
			err := json.Unmarshal([]byte(queryInviteCode.Extra), &extraRegister)
			if err != nil {
				context.JSON(httpcode.InternalError, gin.H{
					"code": httpcode.InternalError,
					"msg":  "服务器内部错误，请稍后重试",
				})
				return
			}
		}
		// 6、注册用户
		newUser, result := database.UserCreateAndSave(request.Username, request.Password, extraRegister.UserType)
		if result.Code != errcode.OK {
			context.JSON(httpcode.InternalError, gin.H{
				"code": result.Code,
				"msg":  "服务器内部错误，请稍后重试",
			})
			return
		}
		// 7、返回信息
		context.JSON(httpcode.OK, gin.H{
			"code": httpcode.OK,
			"msg":  "注册成功",
			"data": gin.H{
				"user_id":     newUser.ID,
				"username":    newUser.Username,
				"nickname":    newUser.Nickname,
				"avatar":      newUser.Avatar,
				"email":       newUser.Email,
				"phone":       newUser.Phone,
				"create_time": newUser.CreateTime,
				"type":        newUser.Type,
				"status":      newUser.Status,
			},
		})
	})

}
