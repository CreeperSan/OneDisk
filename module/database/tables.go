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

const valueUserTypeGuest = 0
const valueUserTypeNormal = 1
const valueUserTypeAdmin = 2

type User struct {
	ID         int64  `gorm:"column:id;"`
	Username   string `gorm:"column:username;"`
	Password   string `gorm:"column:password;"`
	Nickname   string `gorm:"column:nickname;"`
	Avatar     string `gorm:"column:avatar;"`
	Email      string `gorm:"column:email;"`
	Phone      string `gorm:"column:phone;"`
	CreateTime int64  `gorm:"column:create_time;"`
	Type       int    `gorm:"column:type;"` // 用户类型（0：游客，1：普通用户，2：管理员）
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
const columnUserTokenSecretKey = "secret_key"
const columnUserTokenCreateTime = "create_time"
const columnUserTokenExpireTime = "expire_time"

const valueUserTokenPlatformUnknown = 0
const valueUserTokenPlatformBrowser = 1
const valueUserTokenPlatformAndroid = 2
const valueUserTokenPlatformIOS = 3
const valueUserTokenPlatformWindows = 4
const valueUserTokenPlatformMacOS = 5
const valueUserTokenPlatformLinux = 6

type UserToken struct {
	ID          int64  `gorm:"column:id;"`
	UserID      int64  `gorm:"column:user_id;"`
	Token       string `gorm:"column:token;"`
	Platform    int    `gorm:"column:platform;"`
	MachineCode string `gorm:"column:machine_code;"`
	MachineName string `gorm:"column:machine_name;"`
	SecretKey   string `gorm:"column:secret_key;"`
	CreateTime  int64  `gorm:"column:create_time;"`
	ExpireTime  int64  `gorm:"column:expire_time;"`
}
