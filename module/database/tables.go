package database

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////// 用户标识表

const tableUser = "user"
const columnUserID = "id"
const columnUserUsername = "username"
const columnUserPassword = "password"
const columnUserNickname = "nickname"
const columnUserAvatar = "avatar"
const columnUserEmail = "email"
const columnUserPhone = "phone"
const columnUserCreateTime = "create_time"
const columnUserType = "type"
const columnUserStatus = "Status"

const valueUserTypeGuest = 0
const valueUserTypeNormal = 1
const valueUserTypeAdmin = 2

const valueUserStatusActive = 0
const valueUserStatusForbidden = 1

type User struct {
	ID         int64  `gorm:"column:id;"`
	Username   string `gorm:"column:username;"`
	Password   string `gorm:"column:password;"`
	Nickname   string `gorm:"column:nickname;"`
	Avatar     string `gorm:"column:avatar;"`
	Email      string `gorm:"column:email;"`
	Phone      string `gorm:"column:phone;"`
	CreateTime int64  `gorm:"column:create_time;"`
	Type       int    `gorm:"column:type;"`   // 用户类型（0：游客，1：普通用户，2：管理员）
	Status     int    `gorm:"column:status;"` // 用户状态（0：活跃，1：封禁）
}

func (User) TableName() string {
	return tableUser
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////// 用户登录令牌表

const tableUserToken = "user_token"
const columnUserTokenID = "id"
const columnUserTokenUserID = "user_id"
const columnUserTokenToken = "token"
const columnUserTokenPlatform = "platform"
const columnUserTokenMachineCode = "machine_code"
const columnUserTokenMachineName = "machine_name"
const columnUserTokenRefreshToken = "refresh_token"
const columnUserTokenCreateTime = "create_time"
const columnUserTokenTokenExpireTime = "token_expire_time"
const columnUserTokenRefreshTokenExpireTime = "refresh_token_expire_time"
const columnUserTokenLastAccessTime = "last_access_time"
const columnUserTokenLastRefreshTime = "last_refresh_time"

const valueUserTokenPlatformUnknown = 0
const valueUserTokenPlatformBrowser = 1
const valueUserTokenPlatformAndroid = 2
const valueUserTokenPlatformIOS = 3
const valueUserTokenPlatformWindows = 4
const valueUserTokenPlatformMacOS = 5
const valueUserTokenPlatformLinux = 6

type UserToken struct {
	ID                     int64  `gorm:"column:id;"`
	UserID                 int64  `gorm:"column:user_id;"`
	Token                  string `gorm:"column:token;"`
	Platform               int    `gorm:"column:platform;"`
	MachineCode            string `gorm:"column:machine_code;"`
	MachineName            string `gorm:"column:machine_name;"`
	RefreshToken           string `gorm:"column:refresh_token;"`
	TokenExpireTime        int64  `gorm:"column:token_expire_time;"`
	RefreshTokenExpireTime int64  `gorm:"column:refresh_token_expire_time;"`
	CreateTime             int64  `gorm:"column:create_time;"`
	LastAccessTime         int64  `gorm:"column:last_access_time;"`
	LastRefreshTime        int64  `gorm:"column:last_refresh_time;"`
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////// 用户注册邀请码

const tableUserInviteCode = "user_invite_code"
const columnUserInviteCodeID = "id"
const columnUserInviteCodeFromUserID = "from_user_id"
const columnUserInviteCodeExpiredTime = "expired_time"
const columnUserInviteCodeUsage = "usage"
const columnUserInviteCodeCode = "code"

const valueUserInviteCodeUsageRegister = "register"

type UserInviteCode struct {
	ID          int64  `gorm:"column:id;"`
	FromUserID  int64  `gorm:"column:from_user_id;"`
	ExpiredTime int64  `gorm:"column:expired_time;"`
	Usage       string `gorm:"column:usage;"`
	code        string `gorm:"column:code;"`
}
